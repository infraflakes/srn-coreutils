{
  description = "An opinionated CLI wrapper that replaces cryptic flags with self-explanatory, English-like sub-commands.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    srn-coreutils-test = {
      url = "github:nixuris/srn-coreutils?ref=main";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    srn-coreutils-test,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        # Version and download info for stable release
        stableVersion = "3.0.0"; # Change this
        stableDownload = {
          url = "https://github.com/nixuris/srn-coreutils/releases/download/v${stableVersion}/serein_${stableVersion}_linux_amd64.tar.gz";
          sha256 = "vKPHNDIXz1YC5hbLU88i/vkRHKAO019SyKv9vTnw53w=";
        };

        # Source build for test variant
        buildSereinFromSource = {
          src,
          version,
        }:
          pkgs.buildGoModule {
            pname = "serein";
            inherit version src;
            vendorHash = "sha256-+gNaABMs7XZbOFlvLQA5KtnZrBHDWgBtH6W29KMeBU0="; # Update if source changes
            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
            ];
            nativeBuildInputs = [pkgs.installShellFiles];
            postFixup = ''
              installShellCompletion --fish ${src}/completions/serein.fish
              installShellCompletion --zsh ${src}/completions/serein.zsh
              installShellCompletion --bash ${src}/completions/serein.bash
            '';
          };

        # Binary install for stable variant
        buildSereinFromBinary = pkgs.stdenv.mkDerivation {
          pname = "serein";
          version = stableVersion;
          src = builtins.fetchurl {
            inherit (stableDownload) url sha256;
          };
          nativeBuildInputs = [
            pkgs.installShellFiles
            pkgs.patchelf
            pkgs.makeWrapper
          ];
          phases = [
            "installPhase"
            "fixupPhase"
          ];
          installPhase = ''
            mkdir -p $out/bin
            tar -xzf $src
            cp serein $out/bin/serein
            chmod +w $out/bin/serein  # Ensure the binary is writable for patchelf
            chmod +x $out/bin/serein
            installShellCompletion --fish ${./.}/completions/serein.fish
            installShellCompletion --zsh ${./.}/completions/serein.zsh
            installShellCompletion --bash ${./.}/completions/serein.bash
          '';
          fixupPhase = ''
            patchelf --set-interpreter "$(cat $NIX_CC/nix-support/dynamic-linker)" $out/bin/serein
            wrapProgram $out/bin/serein --prefix LD_LIBRARY_PATH : "${pkgs.lib.makeLibraryPath [pkgs.glibc]}"
          '';
        };

        # Clean source with explicit inclusion of .version
        cleanedSource = pkgs.lib.cleanSourceWith {
          src = ./.;
          filter = path: type: let
            baseName = baseNameOf path;
          in
            baseName == ".version" || pkgs.lib.cleanSourceFilter path type;
        };
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint
            cmake
          ];
        };
        packages = rec {
          # Test variant: Build from local source
          test = buildSereinFromSource {
            src = cleanedSource;
            version = let
              versionFile = "${cleanedSource}/.version";
            in
              pkgs.lib.escapeShellArg (
                if builtins.pathExists versionFile
                then builtins.readFile versionFile
                else self.shortRev or "dev"
              );
          };

          # Stable variant: Install from pre-built binary
          stable = buildSereinFromBinary;

          # Default package
          default = test;
        };

        apps = rec {
          default = flake-utils.lib.mkApp {drv = self.packages.${system}.default;};
          test = flake-utils.lib.mkApp {drv = self.packages.${system}.test;};
          stable = flake-utils.lib.mkApp {drv = self.packages.${system}.stable;};
        };
      }
    );
}

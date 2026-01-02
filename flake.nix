{
  description = "Serein is an opinionated CLI suite to streamline many command line work.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        buildSrn = {
          src,
          version,
        }:
          pkgs.buildGoModule {
            pname = "srn";
            inherit version src;
            preBuild = ''
              export CGO_ENABLED=0
            '';
            vendorHash = "sha256-rOxGSxxs8R3355AyQG6+hsbtKOw33tmswSSVyzy9JA8="; # Update if source changes
            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
            ];
            nativeBuildInputs = [pkgs.installShellFiles];
            postInstall = ''
              mv $out/bin/srn-coreutils $out/bin/srn
            '';
            postFixup = ''
              installShellCompletion --fish ${src}/completions/srn.fish
              installShellCompletion --zsh ${src}/completions/srn.zsh
              installShellCompletion --bash ${src}/completions/srn.bash
            '';
          };

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
            goreleaser
          ];
        };

        packages.default = buildSrn {
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

        apps.default = flake-utils.lib.mkApp {drv = self.packages.${system}.default;};
      }
    );
}

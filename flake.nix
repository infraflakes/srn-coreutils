{
  description = "An opinionated CLI wrapper that replaces cryptic flags with self-explanatory, English-like sub-commands.";

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

        buildSerein = {
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

        packages.default = buildSerein {
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

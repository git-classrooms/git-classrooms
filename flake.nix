{
  description = "Development enviroment for git classrooms";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    pre-commit-hooks.url = "github:cachix/git-hooks.nix";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  } @ inputs:
    {
      overlays.default = final: prev: {
        git-classrooms = prev.callPackage ./nix/package.nix {};
      };

      nixosModules = {
        default = import ./nix/module.nix;
        git-classrooms = import ./nix/module.nix;
      };
    }
    // (flake-utils.lib.eachDefaultSystem
      (system: let
        pkgs = nixpkgs.legacyPackages.${system};
        pre-commit-check = inputs.pre-commit-hooks.lib.${system}.run {
          src = ./.;
          hooks = {
            golangci-lint.enable = false;
            gofmt.enable = true;
            gotest.enable = false;
          };
        };
      in {
        packages.default = pkgs.callPackage ./nix/package.nix {};

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Frontend
            nodejs_24
            pnpm

            # Backend and tools
            go
            delve

            yq-go
            docker
            docker-compose
            gnumake
            git
            postgresql
          ];

          shellHook = ''
            echo "Welcome to the Git Classrooms dev shell"

            (cd ./frontend && pnpm install)
            go mod download
            ${pre-commit-check.shellHook}
          '';
        };
      }));
}

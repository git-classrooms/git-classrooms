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
  } @ inputs: (flake-utils.lib.eachDefaultSystem
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
          cd ./frontend && pnpm install

          cd ..
          go mod download
          ${pre-commit-check.shellHook}
        '';
      };
    }));
}

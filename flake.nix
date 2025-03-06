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
      devShells.default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          # Frontend
          nodejs_22
          yarn

          # Backend and tools
          go_1_24
          delve

          yq-go
          docker
          docker-compose
          gnumake
          git
          postgresql
        ];

        shellHook = ''
          cd ./frontend && yarn install

          cd ..
          go mod download
          ${pre-commit-check.shellHook}
        '';
      };
    }));
}

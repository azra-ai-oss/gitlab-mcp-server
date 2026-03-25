{
  description = "GitLab MCP Server - A Go implementation of Model Context Protocol for GitLab";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            gopls
            gotools
            go-tools # staticcheck
            golangci-lint
            pre-commit
          ];

          shellHook = ''
            echo "GitLab MCP Server dev shell"
            echo "Go version: $(go version)"

            # Auto-install pre-commit hooks
            if [ -d .git ] && [ -f .pre-commit-config.yaml ]; then
              pre-commit install > /dev/null 2>&1
              echo "Pre-commit hooks installed"
            fi
          '';
        };

        packages.default = pkgs.buildGo125Module {
          pname = "gitlab-mcp-server";
          version = "0.2.0";
          src = ./.;
          vendorHash = null; # Update after first build
        };
      }
    );
}

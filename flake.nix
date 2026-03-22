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
          ];

          shellHook = ''
            echo "GitLab MCP Server dev shell"
            echo "Go version: $(go version)"
          '';
        };

        packages.default = pkgs.buildGo125Module {
          pname = "gitlab-mcp-server";
          version = "0.1.0";
          src = ./.;
          vendorHash = null; # Update after first build
        };
      }
    );
}

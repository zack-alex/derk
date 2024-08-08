{
  description = "derk";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ ];
        };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = [ pkgs.go pkgs.gomplate ];
        };
        packages.default = pkgs.buildGoModule {
          name = "derk";

          src = ./.;

          buildInputs = [ pkgs.go pkgs.gomplate ];

          vendorHash = "sha256-D7hqPo8HlEqEnF4TagLQXkUVJZq+2An1h3trPpJoD5Q=";
        };
      }
    );
}


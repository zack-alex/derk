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
        packages.default = pkgs.stdenv.mkDerivation {
          name = "derk";

          src = ./.;

          buildInputs = [ pkgs.go ];

          buildPhase = ''
            export GOCACHE=$PWD/temp/cache
            export GOMODCACHE=$PWD/temp/modcache
            ./scripts/build
          '';

          installPhase = ''
            mkdir -p $out/bin
            mv ./dev/derk $out/bin/
          '';
        };
      }
    );
}


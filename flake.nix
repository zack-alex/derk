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

          buildInputs = [ pkgs.go pkgs.gomplate ];

          buildPhase = ''
            export GOCACHE=$PWD/temp/cache
            export GOMODCACHE=$PWD/temp/modcache
            pwd
            ls
            echo FOOO
            ls scripts
            ls scripts/build
            cat scripts/build
            bash ./scripts/build
            echo build ok
          '';

          installPhase = ''
            mkdir -p $out/bin $out/share
            cp ./dev/derk $out/bin/
            cp ./dev/index.html $out/share/
          '';
        };
      }
    );
}


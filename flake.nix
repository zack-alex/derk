{
  description = "A Nix flake for the password-deriver Python package";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ ];
        };
        pythonEnv = pkgs.python3.withPackages (ps: [ ]);
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = [ pythonEnv pkgs.python3Packages.setuptools pkgs.python3Packages.pytest pkgs.python3Packages.pip ];
        };

        packages.default = pkgs.python3Packages.buildPythonPackage {
          pname = "password-deriver";
          version = "0.0.1";
          src = self;
          propagatedBuildInputs = [ pythonEnv ];
        };
      }
    );
}


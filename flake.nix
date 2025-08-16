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
        pkgsWasi32 = pkgs.pkgsCross.wasi32;
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
        packages.derk-web = pkgsWasi32.buildGoModule {
          name = "derk-web";

          src = ./.;
          subPackages = [
            "web"
          ];

          nativeBuildInputs = [
            pkgs.go
            pkgs.gomplate
          ];

          vendorHash = "sha256-6k7HaEK63HmC2YZA9dAL3zz4Gypc2ygPkWCigqNrfeg=";

          preBuild = ''
            export GOOS=js
          '';
          postInstall = ''
            base64 -w 0 $out/bin/web >$out/bin/derk.base64
            cp -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" $out/bin/wasm_exec.js
            cp ./web/index.html.templ $out/bin
            cd $out/bin
            gomplate -f index.html.templ -o index.html
            mv index.html ..
            rm *
            mv ../index.html .
          '';
        };
      }
    );
}


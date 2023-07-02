{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
}:

let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  packages = [
    goEnv

    # https://pkg.go.dev/golang.org/x/tools/gopls
    pkgs.gopls

    # https://pkg.go.dev/github.com/ramya-rao-a/go-outline
    pkgs.go-outline

    # https://github.com/cweill/gotests
    pkgs.gotests

    # https://github.com/go-delve/delve
    pkgs.delve

    # goimports, godoc, etc.
    pkgs.gotools

    # https://github.com/golangci/golangci-lint
    pkgs.golangci-lint

    # https://pkg.go.dev/github.com/josharian/impl
    pkgs.impl

    # https://github.com/nix-community/gomod2nix
    pkgs.gomod2nix
  ];
  shellHook = ''
    ${pkgs.go}/bin/go version
  '';
}

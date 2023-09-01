{
  description = "A Nix-flake-based development environment for Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    nixpkgs,
    gitignore,
    gomod2nix,
  }: let
    goVersion = 20; # Change this to update the whole stack
    overlays = [
      # Overlay our custom Go version
      (final: prev: {go = prev."go_1_${toString goVersion}";})
      # Overlay gomod2nix so that we can use buildGoApplication
      gomod2nix.overlays.default
    ];
    supportedSystems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
    # Provide system-specific attributes
    forAllSystems = f:
      nixpkgs.lib.genAttrs supportedSystems (system:
        f {
          pkgs = import nixpkgs {inherit overlays system;};
        });
  in {
    packages = forAllSystems ({pkgs}: {
      default = pkgs.buildGoApplication {
        pname = "dinosaur";
        version = "0.1";
        src = gitignore.lib.gitignoreSource ./.;
        modules = ./gomod2nix.toml;
      };
    });

    devShells = forAllSystems ({pkgs}: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          # go 1.20 (specified by overlay)
          go

          # https://pkg.go.dev/golang.org/x/tools/gopls
          gopls

          # https://pkg.go.dev/github.com/ramya-rao-a/go-outline
          go-outline

          # https://github.com/cweill/gotests
          gotests

          # https://github.com/go-delve/delve
          delve

          # goimports, godoc, etc.
          gotools

          # https://github.com/golangci/golangci-lint
          golangci-lint

          # https://pkg.go.dev/github.com/josharian/impl
          impl

          # https://github.com/nix-community/gomod2nix
          # This allows us to run gomod2nix from the dev shell
          gomod2nix.packages.${system}.default
        ];
        shellHook = ''
          ${pkgs.go}/bin/go version
        '';
      };
    });
  };
}

{
  description = "A smart workspace navigation script for Niri";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    utils,
  }:
    utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
        niri-smart-workspace = pkgs.callPackage ./package.nix {};
      in {
        packages.default = niri-smart-workspace;
        packages.niri-smart-workspace = niri-smart-workspace;

        devShells.default = pkgs.mkShell {
          inputsFrom = [niri-smart-workspace];

          nativeBuildInputs = with pkgs; [
            go
          ];
        };
      }
    )
    // {
      homeModules.default = ./module.nix;
    };
}

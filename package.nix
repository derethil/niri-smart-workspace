{
  pkgs,
  lib,
  ...
}:
pkgs.buildGoModule {
  pname = "niri-smart-workspace";
  version = "0.1.0";

  src = ./.;

  vendorHash = null;

  meta = {
    description = "Smart workspace navigation for niri compositor that skips empty workspaces";
    mainProgram = "niri-smart-workspace";
    platforms = lib.platforms.linux;
  };
}

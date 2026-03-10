{
  config,
  lib,
  pkgs,
  ...
}: let
  cfg = config.services.niri-smart-workspace;
in {
  options.services.niri-smart-workspace = {
    enable = lib.mkEnableOption "niri-smart-workspace daemon";
    package = lib.mkOption {
      type = lib.types.package;
      default = pkgs.callPackage ./package.nix {};
      description = "The niri-smart-workspace package to use.";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [cfg.package];

    systemd.user.services.niri-smart-workspace = {
      Unit = {
        Description = "Smart workspace navigation for niri";
        After = ["graphical-session.target"];
      };
      Service = {
        ExecStart = "${lib.getExe cfg.package} --daemon";
        Restart = "on-failure";
      };
      Install.WantedBy = ["graphical-session.target"];
    };
  };
}

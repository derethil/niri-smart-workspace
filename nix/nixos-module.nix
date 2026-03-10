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
    environment.systemPackages = [cfg.package];

    systemd.user.services.niri-smart-workspace = {
      description = "Smart workspace navigation for niri";
      after = ["graphical-session.target"];
      partOf = ["graphical-session.target"];

      serviceConfig = {
        ExecStart = "${lib.getExe cfg.package} --daemon";
        Restart = "on-failure";
        Type = "simple";
      };

      wantedBy = ["graphical-session.target"];
    };
  };
}

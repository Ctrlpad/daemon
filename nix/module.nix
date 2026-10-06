self:
{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.ctrlpad-daemon;
in
{
  options.services.ctrlpad-daemon = {
    enable = lib.mkEnableOption "the Ctrlpad daemon";
    device = lib.mkOption {
      type = lib.types.str;
      default = "";
      description = "MAC address of the Ctrlpad";
    };
  };

  config = lib.mkIf cfg.enable {
    hardware.bluetooth.enable = lib.mkDefault true;

    systemd.user.services.ctrlpad-daemon = {
      description = "Ctrlpad Daemon";
      documentation = [ "https://github.com/ctrlpad/daemon" ];
      wantedBy = [ "graphical-session.target" ];
      partOf = [ "graphical-session.target" ];
      after = [ "graphical-session.target" ];

      serviceConfig = {
        ExecStart = "${lib.getExe self.packages.${pkgs.stdenv.hostPlatform.system}.ctrlpad-daemon} -device ${cfg.device};
        Restart = "on-failure";
        RestartSec = 5;
        Type = "simple";
      };
    };
  };
}

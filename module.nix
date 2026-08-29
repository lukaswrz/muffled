self:
{
  config,
  pkgs,
  lib,
  utils,
  ...
}:
let
  cfg = config.services.muffled;
  format = pkgs.formats.toml { };
  inherit (lib) types;
in
{
  options = {
    services.muffled = {
      enable = lib.mkEnableOption "Muffled";

      package = lib.mkPackageOption self.packages.${pkgs.stdenv.hostPlatform.system} "default" { };

      settings = lib.mkOption {
        default = { };
        # TOML does not allow null values, so we use null to omit those fields
        apply = lib.filterAttrsRecursive (_: v: v != null);
        description = "Settings for Muffled.";
        type = types.submodule {
          freeformType = format.type;

          options = {
            user = lib.mkOption {
              description = "The ListenBrainz user.";
              type = types.nonEmptyStr;
            };

            listen = lib.mkOption {
              default = "localhost:8080";
              description = "The address and port to listen on.";
              type = types.str;
            };

            log_level = lib.mkOption {
              default = "info";
              description = "The log level.";
              type = types.enum [
                "debug"
                "info"
                "warn"
                "error"
              ];
            };

            interval = lib.mkOption {
              default = 120;
              description = "The interval for polling ListenBrainz.";
              type = types.ints.positive;
            };

            listenbrainz_base_url = lib.mkOption {
              default = "https://api.listenbrainz.org/1";
              description = "The ListenBrainz base URL.";
              type = types.str;
            };

            widget_path = lib.mkOption {
              default = null;
              description = "The path to a custom widget.";
              type = types.nullOr types.path;
            };
          };
        };
      };
    };
  };

  config = lib.mkIf cfg.enable {
    systemd = {
      services.muffled = {
        after = [ "network.target" ];
        description = "Muffled";
        wantedBy = [ "multi-user.target" ];
        serviceConfig = {
          ExecStart = utils.escapeSystemdExecArgs [
            (lib.getExe cfg.package)
            "--config"
            (format.generate "config.toml" cfg.settings)
          ];

          DynamicUser = true;

          # Hardening
          NoNewPrivileges = true;
          PrivateDevices = true;
          DevicePolicy = "closed";
          CapabilityBoundingSet = "";
          LockPersonality = true;
          MemoryDenyWriteExecute = true;
          PrivateUsers = true;
          ProtectClock = true;
          ProtectControlGroups = true;
          ProtectHome = true;
          ProtectHostname = true;
          ProtectKernelLogs = true;
          ProtectKernelModules = true;
          ProtectKernelTunables = true;
          ProtectProc = "invisible";
          ProcSubset = "pid";
          ProtectSystem = "strict";
          RestrictAddressFamilies = [
            "AF_INET"
            "AF_INET6"
          ];
          RestrictNamespaces = true;
          RestrictRealtime = true;
          RestrictSUIDSGID = true;
          SystemCallArchitectures = "native";
        };
      };
    };
  };
}

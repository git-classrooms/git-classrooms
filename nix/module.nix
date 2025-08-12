{
  pkgs,
  config,
  lib,
  ...
}: let
  cfg = config.services.git-classrooms;
  nonEmptyStr = lib.types.strMatching ".+";
  urlStr = lib.types.strMatching "^https?://.+";
  durationStr = lib.types.strMatching "^[0-9]+(s|m|h|d)$";
in {
  options.services.git-classrooms = {
    enable = lib.mkEnableOption "Starts the Git Classrooms service as a systemd unit.";

    package = lib.mkOption {
      type = lib.types.package;
      default = pkgs.git-classrooms;
      description = "Package providing the git-classrooms binary to run.";
    };

    dataDir = lib.mkOption {
      type = lib.types.str;
      default = "/var/lib/git-classrooms";
      description = ''
        The directory where git-classrooms stores its stateful data.
      '';
    };

    settings = lib.mkOption {
      type = lib.types.submodule {
        options = {
          publicUrl = lib.mkOption {
            type = urlStr;
            description = ''
              Public base URL of the application including the scheme.
              Used for link generation and OAuth redirect URLs.
            '';
            example = "https://classrooms.example.edu";
          };

          port = lib.mkOption {
            type = lib.types.port;
            default = 3000;
            description = "HTTP listening port of the web server.";
            example = 8080;
          };

          frontendPath = lib.mkOption {
            type = lib.types.str;
            default = "./public";
            description = ''
              Path to the frontend static files directory.
              Usually not needed when using the packaged version.
            '';
          };

          trustedProxies = lib.mkOption {
            type = lib.types.listOf lib.types.str;
            default = [];
            description = ''
              List of trusted proxy IP addresses or CIDR ranges.
              Used for X-Forwarded-For header validation.
            '';
            example = ["127.0.0.1" "192.168.1.0/24"];
          };

          auth = lib.mkOption {
            type = lib.types.submodule {
              options = {
                clientId = lib.mkOption {
                  type = nonEmptyStr;
                  description = "OAuth/OIDC client ID registered at the identity provider.";
                  example = "git-classrooms-web";
                };

                clientSecret = lib.mkOption {
                  type = nonEmptyStr;
                  description = "OAuth/OIDC client secret registered at the identity provider.";
                };

                redirectEndpoint = lib.mkOption {
                  type = lib.types.str;
                  default = "/api/v1/auth/gitlab/callback";
                  description = ''
                    OAuth callback endpoint path (relative to publicUrl).
                  '';
                };

                authUrl = lib.mkOption {
                  type = lib.types.nullOr urlStr;
                  default = null;
                  description = ''
                    OAuth authorization URL. If null, defaults to $GITLAB_URL/oauth/authorize.
                  '';
                  example = "https://gitlab.example.com/oauth/authorize";
                };

                tokenUrl = lib.mkOption {
                  type = lib.types.nullOr urlStr;
                  default = null;
                  description = ''
                    OAuth token URL. If null, defaults to $GITLAB_URL/oauth/token.
                  '';
                  example = "https://gitlab.example.com/oauth/token";
                };

                scopes = lib.mkOption {
                  type = lib.types.listOf nonEmptyStr;
                  default = ["api"];
                  description = ''
                    List of OAuth/OIDC scopes requested during login.
                  '';
                  example = ["api" "read_user"];
                };
              };
            };
            description = "OAuth authentication configuration.";
          };

          gitlab = lib.mkOption {
            type = lib.types.submodule {
              options = {
                url = lib.mkOption {
                  type = urlStr;
                  description = "Base URL of the GitLab instance (e.g., https://gitlab.example.com).";
                  example = "https://gitlab.example.com";
                };
                syncInterval = lib.mkOption {
                  type = durationStr;
                  default = "5m";
                  description = ''
                    Interval for synchronization runs with GitLab.
                    Duration string: "<number><unit>", units: s, m, h, d (e.g., "30s", "5m", "1h").
                    Consider GitLab API rate limits when setting this value.
                  '';
                  example = "1m";
                };
              };
            };
            description = "Configuration for connecting to the GitLab instance.";
          };

          smtp = lib.mkOption {
            type = lib.types.submodule {
              options = {
                host = lib.mkOption {
                  type = nonEmptyStr;
                  description = "SMTP server hostname for outgoing emails.";
                  example = "smtp.example.com";
                };
                port = lib.mkOption {
                  type = lib.types.port;
                  default = 587;
                  description = "SMTP port (e.g., 465 for SMTPS or 587 for STARTTLS).";
                  example = 587;
                };
                user = lib.mkOption {
                  type = lib.types.str;
                  default = "";
                  description = "SMTP username for authentication (empty if not required).";
                  example = "mailer@classrooms.example.edu";
                };
                password = lib.mkOption {
                  type = lib.types.str;
                  default = "";
                  description = "SMTP password for authentication (empty if not required).";
                };
              };
            };
            description = "SMTP configuration for sending emails.";
          };

          database = lib.mkOption {
            type = lib.types.submodule {
              options = {
                postgres = lib.mkOption {
                  type = lib.types.submodule {
                    options = {
                      host = lib.mkOption {
                        type = nonEmptyStr;
                        description = "Hostname or service name of the PostgreSQL server.";
                        example = "localhost";
                      };
                      port = lib.mkOption {
                        type = lib.types.port;
                        default = 5432;
                        description = "PostgreSQL server port.";
                      };
                      user = lib.mkOption {
                        type = nonEmptyStr;
                        default = "postgres";
                        description = "PostgreSQL username.";
                      };
                      password = lib.mkOption {
                        type = nonEmptyStr;
                        description = "PostgreSQL password.";
                      };
                      db = lib.mkOption {
                        type = nonEmptyStr;
                        default = "git-classrooms";
                        description = "Name of the PostgreSQL database.";
                      };
                    };
                  };
                  description = "PostgreSQL database connection settings.";
                };
              };
            };
            description = "Database configuration.";
          };
        };
      };
      description = "Full configuration for the Git Classrooms service.";
    };

    environmentFile = lib.mkOption {
      type = lib.types.nullOr lib.types.path;
      default = null;
      description = ''
        File path to be sourced via systemd `EnvironmentFile=`.
        If set, secret environment variables (e.g. passwords, client secrets)
        should be provided via this file instead of the Nix configuration.
        Variables from this file will override settings from the Nix config.
      '';
      example = "/run/secrets/git-classrooms.env";
    };
  };

  config = lib.mkIf cfg.enable {
    users.users.git-classrooms = {
      group = "git-classrooms";
      home = cfg.dataDir;
      createHome = true;
      isSystemUser = true;
    };

    users.groups.git-classrooms = {};

    systemd.services.git-classrooms = {
      description = "Git Classrooms Service";
      wantedBy = ["multi-user.target"];
      after = ["network-online.target" "postgresql.service"];
      wants = ["network-online.target"];

      serviceConfig = {
        Type = "simple";
        ExecStart = "${cfg.package}/bin/git-classrooms";
        User = "git-classrooms";
        Group = "git-classrooms";
        WorkingDirectory = cfg.dataDir;
        Restart = "always";
        RestartSec = 5;

        # Security hardening
        NoNewPrivileges = true;
        PrivateTmp = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        ReadWritePaths = [cfg.dataDir];

        # Load secrets from file if provided
      } // lib.optionalAttrs (cfg.environmentFile != null) {
        EnvironmentFile = cfg.environmentFile;
      };

      environment = lib.mkMerge [
        # Always set these base configuration values
        {
          PUBLIC_URL = cfg.settings.publicUrl;
          PORT = toString cfg.settings.port;
          FRONTEND_PATH = cfg.settings.frontendPath;

          POSTGRES_HOST = cfg.settings.database.postgres.host;
          POSTGRES_PORT = toString cfg.settings.database.postgres.port;
          POSTGRES_USER = cfg.settings.database.postgres.user;
          POSTGRES_DB = cfg.settings.database.postgres.db;

          GITLAB_URL = cfg.settings.gitlab.url;
          GITLAB_SYNC_INTERVAL = cfg.settings.gitlab.syncInterval;

          SMTP_HOST = cfg.settings.smtp.host;
          SMTP_PORT = toString cfg.settings.smtp.port;

          AUTH_REDIRECT_ENDPOINT = cfg.settings.auth.redirectEndpoint;
          AUTH_SCOPES = lib.concatStringsSep "," cfg.settings.auth.scopes;
        }

        # Only set sensitive values if environmentFile is not used
        (lib.mkIf (cfg.environmentFile == null) {
          POSTGRES_PASSWORD = cfg.settings.database.postgres.password;
          SMTP_USER = cfg.settings.smtp.user;
          SMTP_PASSWORD = cfg.settings.smtp.password;
          AUTH_CLIENT_ID = cfg.settings.auth.clientId;
          AUTH_CLIENT_SECRET = cfg.settings.auth.clientSecret;
        })

        # Set optional values if provided
        (lib.mkIf (cfg.settings.trustedProxies != []) {
          TRUSTED_PROXIES = lib.concatStringsSep "," cfg.settings.trustedProxies;
        })

        (lib.mkIf (cfg.settings.auth.authUrl != null) {
          AUTH_AUTH_URL = cfg.settings.auth.authUrl;
        })

        (lib.mkIf (cfg.settings.auth.tokenUrl != null) {
          AUTH_TOKEN_URL = cfg.settings.auth.tokenUrl;
        })
      ];
    };
  };
}

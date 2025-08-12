{
  stdenv,
  fetchYarnDeps,
  yarnConfigHook,
  yarnBuildHook,
  yarnInstallHook,
  nodejs,
  buildGoModule,
}: let
  version = "1.0.0";
  frontend = stdenv.mkDerivation (finalAttrs: {
    inherit version;
    pname = "git-classrooms-frontend";
    src = ../frontend;

    yarnOfflineCache = fetchYarnDeps {
      yarnLock = finalAttrs.src + "/yarn.lock";
      hash = "sha256-XwU1/w3vOmRq6rV42pHmX+6LX7mUyof5hJYPw4XMv6I=";
    };

    nativeBuildInputs = [
      yarnConfigHook
      yarnBuildHook
      yarnInstallHook
      nodejs
    ];

    installPhase = ''
      mkdir -p $out
      cp -r dist/* $out/
    '';
  });
in
  buildGoModule (finalAttrs: {
    inherit version;
    pname = "git-classrooms";
    src = ./..;

    ldflags = ["-s" "-w" "-X main.version=${finalAttrs.version}"];

    preBuild = ''
      go generate
      cp -r ${frontend} frontend/dist
    '';

    doCheck = false;
    vendorHash = "sha256-d6TPVLJHf/fGOLZpJ0iYEbwDLWyd9adMFjSI9lzMvV4=";
  })

{
  stdenv,
  fetchPnpmDeps,
  pnpmConfigHook,
  pnpm,
  nodejs,
  buildGoModule,
}: let
  version = "1.0.0";
  frontend = stdenv.mkDerivation (finalAttrs: {
    inherit version;
    pname = "git-classrooms-frontend";
    src = ../frontend;

    pnpmDeps = fetchPnpmDeps {
      inherit (finalAttrs) pname version src;
      hash = "sha256-o2KyNl1y4MTG18k+Wv+ksUzRQ7HS+aUfsZ3P63Aq9a0=";
      fetcherVersion = 3;
    };

    nativeBuildInputs = [
      pnpmConfigHook
      pnpm
      nodejs
    ];

    buildPhase = ''
      runHook preBuild
      pnpm build
      runHook postBuild
    '';

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
      go generate ./...
      cp -r ${frontend} frontend/dist
    '';

    postInstall = ''
      mv $out/bin/gitlab-classroom $out/bin/git-classrooms
    '';

    proxyVendor = true;
    doCheck = false;
    vendorHash = "sha256-sy21GRmgKG2gEhXjoysyhPvwVnCz6/6oTvdl5Mxtte8=";
  })

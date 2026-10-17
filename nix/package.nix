{
  lib,
  buildGoModule,
  version ? "0.0.0",
  commit ? "none",
  date ? "unknown",
}:

buildGoModule {
  pname = "ctrlpad-daemon";
  inherit version;

  src = lib.fileset.toSource {
    root = ../.;
    fileset = lib.fileset.unions [
      ../go.mod
      ../go.sum
      ../main.go
      ../cmd
      ../internal
    ];
  };

  vendorHash = "sha256-YuJOts724x2eSJcT/Wi/mIdeqdvBvxz8wzNV+ZsA/YY=";

  env.CGO_ENABLED = 0;

  postInstall = ''
    mv $out/bin/daemon $out/bin/ctrlpad-daemon
  '';

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
    "-X main.commit=${commit}"
    "-X main.date=${date}"
  ];

  meta = {
    description = "Background service that listens for and executes Ctrlpad button actions";
    homepage = "https://github.com/ctrlpad/daemon";
    license = lib.licenses.mit;
    mainProgram = "ctrlpad-daemon";
    platforms = lib.platforms.linux;
  };
}

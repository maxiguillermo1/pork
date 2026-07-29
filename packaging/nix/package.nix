{
  lib,
  buildGoModule,
  fetchFromGitHub,
  version ? "0.2.0",
  source ? fetchFromGitHub {
    owner = "maxiguillermo1";
    repo = "pork";
    tag = "v${version}";
    hash = "sha256-Z6F9r/Hfi4zsi0VH9Z7dEn1nA8B+aVKYqiq72J4m5BU=";
  },
}:

buildGoModule rec {
  pname = "pork";
  inherit version;

  src = source;

  vendorHash = "sha256-Vk3lmPUDvuhOzha8GlH0anRLgVhhyFjz9y328T2gycs=";

  subPackages = [ "cmd/pork" ];

  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];

  # Nix builds with HOME=/homeless-shelter. A few integration-style tests create
  # an isolated pork config, which derives its default download directory from
  # HOME, so give the check phase a writable private home instead.
  preCheck = ''
    export HOME="$TMPDIR/home"
    mkdir -p "$HOME"
  '';

  meta = {
    description = "Terminal torrent search and download client";
    homepage = "https://github.com/maxiguillermo1/pork";
    license = lib.licenses.mit;
    mainProgram = "pork";
    platforms = lib.platforms.unix;
  };
}

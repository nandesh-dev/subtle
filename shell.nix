{
  pkgs ? import <nixpkgs> { },
}:
let
  unstableTarball = fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz";
  unstable = import unstableTarball { };
in
pkgs.mkShell {
  buildInputs = [
    pkgs.ffmpeg_7-headless
    pkgs.tesseract
    pkgs.leptonica
    pkgs.protobuf
    pkgs.protoc-gen-go
    pkgs.protoc-gen-connect-go
    pkgs.graphviz
    pkgs.nodejs_22
    pkgs.pnpm
    unstable.go
  ];

  shellHook = ''
    go install github.com/bufbuild/buf/cmd/buf@v1.46.0
    [ -n "$(go env GOPATH)" ] && export PATH="$(go env GOPATH)/bin:$PATH"

    export CONFIG_DIRECTORY="config"
  '';
}

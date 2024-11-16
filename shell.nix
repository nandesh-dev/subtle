{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShell {
  buildInputs = (
    with pkgs;
    [
      go
      ffmpeg_7-headless
      tesseract
      leptonica
      protobuf
      protoc-gen-go
      protoc-gen-connect-go
      graphviz
      nodejs_22
      pnpm
    ]
  );

  shellHook = ''
    go install github.com/bufbuild/buf/cmd/buf@v1.46.0
    [ -n "$(go env GOPATH)" ] && export PATH="$(go env GOPATH)/bin:$PATH"

    export CONFIG_DIRECTORY="config"
  '';
}

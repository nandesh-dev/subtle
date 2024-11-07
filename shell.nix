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
    go install github.com/bufbuild/buf/cmd/buf@latest
    [ -n "$(go env GOPATH)" ] && export PATH="$(go env GOPATH)/bin:$PATH"

    export CONFIG_DIRECTORY="config"
  '';
}

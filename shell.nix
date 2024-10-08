{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShell {
  nativeBuildInputs = with pkgs.buildPackages; [
    go
    ffmpeg_7-headless
    tesseract
    leptonica
    protobuf
    protoc-gen-go-grpc
    protoc-gen-go
    graphviz
  ];
}

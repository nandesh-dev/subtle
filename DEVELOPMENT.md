# Development

## Setup

1. Install dependencies listed in [README.md](https://github.com/nandesh-dev/subtle/blob/main/README.md) or [shell.nix](https://github.com/nandesh-dev/subtle/blob/main/shell.nix)

2. Generate protobuf files 

```bash
make proto
```

3. Run go 

```bash
go run ./cmd/subtle
```

## Docker Image

### Volumes

|Path|Use|
|---|---|
|`/media`|Directory where all your media content is stored|
|`/config`|Directory where all the subtle data will be stored|

## References

### Golang Packages
[gosseract](https://pkg.go.dev/github.com/otiai10/gosseract/v2@v2.4.1)
[ffmpeg\_go](https://pkg.go.dev/github.com/u2takey/ffmpeg-go)
[grpc](https://grpc.io/docs/languages/go/)
[yaml](https://pkg.go.dev/gopkg.in/yaml.v3)

### File Formats
[Presentation Graphic Stream](https://blog.thescorpius.com/index.php/2017/07/15/presentation-graphic-stream-sup-files-bluray-subtitle-format/)
[Sub Station Alpha v4.00+](http://www.tcax.org/docs/ass-specs.htm)
[SubRip](https://docs.fileformat.com/video/srt/)

### Other
[tesseract](https://github.com/tesseract-ocr/tesseract)

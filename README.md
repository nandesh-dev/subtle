# Subtle

Your tool to manage all the subtitles in your personal media library. Subtle helps you extract and organize subtitle files for all your media content.

## Building

### Debain [ Untested ]

1. Install [Golang](https://go.dev/doc/install)

2. Install dependencies

```bash
apt install tesseract-ocr libtesseract-dev libleptonica-dev ffmpeg
```

3. Build the package

```bash
go build ./cmd/subtle
```

### Docker

```bash
docker run -v /path/to/config:/config -v /path/to/media:/media nandeshdev/subtle
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

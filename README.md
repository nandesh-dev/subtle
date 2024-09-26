# Subtle

Your tool to manage all the subtitles in your personal media library. Subtle helps you extract and organize subtitle files for all your media content.

## Building

### Debain [ Untested ]

1. Install [Golang](https://go.dev/doc/install)

2. Install dependencies

```bash
    apt install tesseract-ocr libtesseract-dev libleptonica-dev
```

3. Build the package

```bash
go build ./cmd/tui
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

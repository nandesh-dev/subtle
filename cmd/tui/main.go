package main

import (
	"log"

	"github.com/nandesh-dev/subtle/internal/filemanager"
	"github.com/nandesh-dev/subtle/internal/subtitle/srt"
)

func main() {
	dir, _ := filemanager.ReadDirectory(".")

	stats, _ := dir.Videos[0].Stats()

	subtitle, err := dir.Videos[0].ExtractSubtitle(stats.Streams[0])

	if err != nil {
		log.Fatal(err)
	}

	srt.EncodeSRTSubtitles(*subtitle)
}

package main

import (
	"fmt"
	"log"

	"github.com/nandesh-dev/subtle/internal/filemanager"
	"github.com/nandesh-dev/subtle/internal/subtitle/parser"
)

func main() {
	dir, _ := filemanager.ReadDirectory(".")

	stats, _ := dir.Videos[0].Stats()

	subtitle, err := parser.ParseRawSubtitleStream(stats.Streams[0])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Subtitle: %v", subtitle)
}

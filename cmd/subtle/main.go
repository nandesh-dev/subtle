package main

import (
	"fmt"

	"github.com/nandesh-dev/subtle/internal/decoder"
	"github.com/nandesh-dev/subtle/internal/filemanager"
)

func main() {
	dir, _ := filemanager.ReadDirectory("./media")

	stats, _ := dir.Videos[0].Stats()

	stream, _, _ := decoder.DecodeRawSubtitleStream(stats.RawStreams[0])

	for _, seg := range stream.Segments {
		fmt.Println(seg.Text)
	}
}

package main

import (
	"fmt"

	"github.com/nandesh-dev/subtle/internal/decoder"
	"github.com/nandesh-dev/subtle/internal/filemanager"
)

func main() {
	dir, _ := filemanager.ReadDirectory("./media")

	stats, _ := dir.Videos[0].Stats()

	fmt.Println(decoder.DecodeRawSubtitleStream(stats.RawStreams[1]))
}

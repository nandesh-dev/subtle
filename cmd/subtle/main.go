package main

import (
	"os"

	"github.com/nandesh-dev/subtle/internal/ass"
	"github.com/nandesh-dev/subtle/internal/filemanager"
	"github.com/nandesh-dev/subtle/internal/srt"
	"github.com/nandesh-dev/subtle/internal/subtitle"
)

func main() {
	dir, _ := filemanager.ReadDirectory("./media")
	videos, _ := dir.VideoFiles()

	for _, v := range videos {
		rawStreams, _ := subtitle.ExtractRawStreams(&v)

		for _, rawStream := range rawStreams {
			if rawStream.Format() == subtitle.ASS {
				assStream, _, _ := ass.DecodeSubtitle(&rawStream)

				srtStream := srt.NewStream()

				for _, seg := range assStream.Segments() {
					srtSegment := srt.NewSegment()
					srtSegment.SetStart(seg.Start())
					srtSegment.SetEnd(seg.End())
					srtSegment.SetText(seg.Text())

					srtStream.AddSegment(*srtSegment)
				}

				file, _ := os.Create("subtitle.srt")
				defer file.Close()

				file.WriteString(srt.EncodeSubtitle(*srtStream))
				break
			}
		}
	}
}

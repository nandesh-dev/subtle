package writer

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

type Writer struct {
	file *os.File
}

func NewWriter(file *os.File) *Writer {
	return &Writer{
		file: file,
	}
}

func (w *Writer) Write(sub *subtitle.Subtitle) error {
	writer := bufio.NewWriter(w.file)

	for i, cue := range sub.Cues {
		indexLine := fmt.Sprintf("%v\n", i)
		if _, err := writer.WriteString(indexLine); err != nil {
			return err
		}

		timestampLine := fmt.Sprintf("%v --> %v\n", formatTimestamp(cue.Timestamp.Start), formatTimestamp(cue.Timestamp.End))
		if _, err := writer.WriteString(timestampLine); err != nil {
			return err
		}

		contentLine := ""
		for i, contentSegment := range cue.Content {
			contentLine = contentSegment.Text
			if i != len(cue.Content)-1 {
				contentLine += " "
			}
		}
		contentLine += "\n\n"

		if _, err := writer.WriteString(contentLine); err != nil {
			return err
		}
	}

	return writer.Flush()
}

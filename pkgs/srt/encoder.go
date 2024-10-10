package srt

import (
	"fmt"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

func EncodeSubtitle(sub subtitle.TextSubtitle) string {
	segments := sub.Segments()

	output := ""

	for i, segment := range segments {
		output += fmt.Sprintln(i)
		output += fmt.Sprintf("%v --> %v\n", formatTimestamp(segment.Start()), formatTimestamp(segment.End()))
		output += fmt.Sprintf("%v\n\n", segment.Text())
	}

	return output
}

func formatTimestamp(t time.Duration) string {
	hours := int(t.Hours())
	minutes := int(t.Minutes()) % 60
	seconds := int(t.Seconds()) % 60
	milliseconds := int(t.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}

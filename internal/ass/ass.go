package ass

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/warnings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type SectionType int

const (
	Info SectionType = iota
	Styles
	Events
)

type LineType int

const (
	SynchPoint LineType = iota
	Timer
	Dialogue
	Format
)

func extractLineTypePrefix(line string) (LineType, string, error) {
	if strings.HasPrefix(line, "Synch Point: ") {
		return SynchPoint, strings.TrimPrefix(line, "Synch Point: "), nil
	} else if strings.HasPrefix(line, "Timer: ") {
		return Timer, strings.TrimPrefix(line, "Timer: "), nil
	} else if strings.HasPrefix(line, "Dialogue: ") {
		return Dialogue, strings.TrimPrefix(line, "Dialogue: "), nil
	} else if strings.HasPrefix(line, "Format: ") {
		return Format, strings.TrimPrefix(line, "Format: "), nil
	}

	return Dialogue, line, fmt.Errorf("Unrecognized line type prefix: %v", line)
}

func parseTime(t string, multipler int, offset time.Duration) (time.Duration, error) {
	pt := strings.Split(t, ":")

	if len(pt) < 3 {
		return time.Duration(0), fmt.Errorf("Invalid timestamp: %v", t)
	}

	hr, err := strconv.Atoi(pt[0])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid hr: %v", pt[0])
	}

	min, err := strconv.Atoi(pt[1])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid min: %v", pt[1])
	}

	sPt := strings.Split(pt[2], ".")

	s, err := strconv.Atoi(sPt[0])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid s: %v", pt[1])
	}

	hs, err := strconv.Atoi(sPt[1])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid hundreth second: %v", pt[1])
	}

	tm := time.Duration((((((hr*60)+min)*60)+s)*1000_000_000+hs*10_000_000)*multipler) + offset

	return tm, nil
}

func ExtractText(str string) string {
	insideBracket := false
	txt := ""
	for _, c := range str {
		if insideBracket {
			if c == '}' {
				insideBracket = false
			}
		} else {
			if c == '{' {
				insideBracket = true
			} else {
				txt += string(c)
			}
		}
	}

	return strings.ReplaceAll(txt, "\\N", "\n")
}

func DecodeSubtitle(rawStream subtitle.RawStream) (*subtitle.Stream, error, *warnings.WarningList) {
	warningList := warnings.NewWarningList()

	var subtitleBuf, errorBuf bytes.Buffer

	ffmpeg.LogCompiledCommand = false
	err := ffmpeg.Input(rawStream.VideoFilePath).
		Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", rawStream.Index), "f": "ass"}).
		WithOutput(&subtitleBuf).
		WithErrorOutput(&errorBuf).
		Run()

	if err != nil {
		return nil, fmt.Errorf("Error extracting subtitles: %v %v", err, errorBuf.String()), warningList
	}

	reader := NewReader(subtitleBuf.Bytes())

	currentFormat := make([]string, 0)
	timeMultiplier := 1
	timeOffset := time.Duration(0)

	stream := subtitle.NewStream(rawStream.Language)

	for !reader.ReachedEnd() {
		line, _ := reader.Advance()

		lT, suffix, err := extractLineTypePrefix(line)

		if err != nil {
			warningList.AddWarning(fmt.Errorf("Error extracting line type prefix: %v", err))
			continue
		}

		switch lT {
		case Format:
			currentFormat = strings.Split(suffix, ",")
		case Dialogue:
			segment := subtitle.Segment{}

			parts := strings.SplitN(suffix, ",", len(currentFormat))
			for i, partName := range currentFormat {
				switch partName {
				case "Start":
					start, err := parseTime(parts[i], timeMultiplier, timeOffset)
					if err != nil {
						warningList.AddWarning(fmt.Errorf("Error parsing start timestamp: %v; %v", err, line))
					} else {
						segment.Start = &start
					}
				case "End":
					end, err := parseTime(parts[i], timeMultiplier, timeOffset)
					if err != nil {
						warningList.AddWarning(fmt.Errorf("Error parsing end timestamp: %v; %v", err, line))
					} else {
						segment.End = &end
					}
				case "Text":
					segment.Text = ExtractText(parts[i])
				}
			}

			stream.AddSegment(segment)
		case Timer:
			multiplier, err := strconv.ParseFloat(suffix, 32)
			if err != nil {
				warningList.AddWarning(fmt.Errorf("Error parsing timer: %v; %v", err, line))
			} else {
				if multiplier != 0.0 {
					timeMultiplier = int(multiplier)
				}
			}
		case SynchPoint:
			synchPoint, err := parseTime(suffix, 1, 0)
			if err != nil {
				warningList.AddWarning(fmt.Errorf("Error parsing synch point: %v; %v", err, line))
			} else {
				timeOffset = synchPoint
			}
		}
	}

	return &stream, nil, warningList
}

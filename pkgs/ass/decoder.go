package ass

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/warning"
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

func DecodeSubtitle(rawStream filemanager.RawStream) (*subtitle.TextSubtitle, warning.WarningList, error) {
	warnings := warning.NewWarningList()

	var subtitleBuf, errorBuf bytes.Buffer

	ffmpeg.LogCompiledCommand = false
	err := ffmpeg.Input(rawStream.Filepath()).
		Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", rawStream.Index()), "f": "ass"}).
		WithOutput(&subtitleBuf).
		WithErrorOutput(&errorBuf).
		Run()

	if err != nil {
		return nil, *warnings, fmt.Errorf("Error extracting subtitles: %v %v", err, errorBuf.String())
	}

	reader := NewReader(subtitleBuf.Bytes())

	currentFormat := make([]string, 0)
	timeMultiplier := 1
	timeOffset := time.Duration(0)

	sub := subtitle.NewTextSubtitle()

	for !reader.ReachedEnd() {
		line, _ := reader.Advance()

		lT, suffix, err := extractLineTypePrefix(line)

		if err != nil {
			warnings.AddWarning(fmt.Errorf("Error extracting line type prefix: %v", err))
			continue
		}

		switch lT {
		case Format:
			currentFormat = make([]string, 0)
			for _, pt := range strings.Split(suffix, ",") {
				currentFormat = append(currentFormat, strings.TrimSpace(pt))
			}
		case Dialogue:
			start, end, text := time.Second*0, time.Second*0, ""

			parts := strings.SplitN(suffix, ",", len(currentFormat))

			for i, partName := range currentFormat {
				switch partName {
				case "Start":
					st, err := parseTime(parts[i], timeMultiplier, timeOffset)
					if err != nil {
						warnings.AddWarning(fmt.Errorf("Error parsing start timestamp: %v; %v", err, line))
					} else {
						start = st
					}
				case "End":
					ed, err := parseTime(parts[i], timeMultiplier, timeOffset)
					if err != nil {
						warnings.AddWarning(fmt.Errorf("Error parsing end timestamp: %v; %v", err, line))
					} else {
						end = ed
					}
				case "Text":
					text = strings.TrimSpace(extractText(parts[i]))
				}
			}

			if text != "" {
				segment := subtitle.NewTextSegment(start, end, text)

				if previousSegment, err := sub.PreviousSegment(); err == nil && segment.Text() == previousSegment.Text() {
					combinedSegment := subtitle.NewTextSegment(previousSegment.Start(), end, text)

					sub.UpdatePreviousSegment(*combinedSegment)
				} else {
					sub.AddSegment(*segment)
				}
			}

		case Timer:
			multiplier, err := strconv.ParseFloat(suffix, 32)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error parsing timer: %v; %v", err, line))
			} else {
				if multiplier != 0.0 {
					timeMultiplier = int(multiplier)
				}
			}
		case SynchPoint:
			synchPoint, err := parseTime(suffix, 1, 0)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error parsing synch point: %v; %v", err, line))
			} else {
				timeOffset = synchPoint
			}
		}
	}

	return sub, *warnings, nil
}

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

func extractText(str string) string {
	insideBracket := false

	bracketStartIndex := 0
	drawingMode := false

	txt := ""
	for i, c := range str {
		if insideBracket {
			if c == '}' {
				insideBracket = false

				if strings.Contains(str[bracketStartIndex:i+1], "\\p1") {
					drawingMode = true
				}
			}
		} else {
			if c == '{' {
				insideBracket = true

				bracketStartIndex = i
				drawingMode = false
			} else {
				if !drawingMode {
					txt += string(c)
				}
			}
		}
	}

	txt = strings.ReplaceAll(txt, "\\N", "\n")
	txt = strings.ReplaceAll(txt, "\\h", "")

	return txt
}

package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (_ *Parser) Parse(data []byte) (*subtitle.Subtitle, error) {
	reader := LineReaderFromBytes(data)
	sub := subtitle.Subtitle{}

	headerConfig := HeaderConfig{
		synchPoint: time.Duration(0),
		timer:      1,
	}

	currentSection := InfoSection
	currentFormat := make([]FieldType, 0)

	for !reader.AtEnd() {
		line, err := reader.Advance()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			switch strings.ToLower(line[1 : len(line)-1]) {
			case "script info":
				currentSection = InfoSection
				continue
			case "events":
				currentSection = EventsSection
				continue
			case "v4 styles", "v4+ styles", "v4 styles+":
				currentSection = StylesSection
				continue
			}
		}

		lineType, lineData, err := parseLineType(line)
		if err != nil {
			continue
		}

		switch currentSection {
		case InfoSection:
			switch lineType {
			case TimerLine:
				multiplier, err := strconv.ParseFloat(lineData, 32)
				if err != nil {
					continue
				}
				if multiplier != 0.0 {
					headerConfig.timer = multiplier
				}
			case SynchPointLine:
				synchPoint, err := parseTimestamp(lineData)
				if err != nil {
					continue
				}
				headerConfig.synchPoint = synchPoint
			}
		case EventsSection:
			switch lineType {
			case FormatLine:
				currentFormat = make([]FieldType, 0)

				for _, fieldTypeString := range strings.Split(lineData, ",") {
					fieldType := parseFieldType(fieldTypeString)
					currentFormat = append(currentFormat, fieldType)
				}
			case DialogueLine:
				cue := subtitle.Cue{}

				fields := strings.SplitN(lineData, ",", len(currentFormat))

				for i, field := range fields {
					if i >= len(currentFormat) {
						break
					}

					switch currentFormat[i] {
					case StartField:
						timestamp, err := parseTimestamp(field)
						if err != nil {
							break
						}

						cue.Timestamp.Start = time.Duration(float64(timestamp)*headerConfig.timer) + headerConfig.synchPoint
					case EndField:
						timestamp, err := parseTimestamp(field)
						if err != nil {
							break
						}

						cue.Timestamp.End = time.Duration(float64(timestamp)*headerConfig.timer) + headerConfig.synchPoint
					case TextField:
            text := field
						//TODO Proper extraction of text and parsing of styles
            text = strings.ReplaceAll(text, "\\N", "\n")
            text = strings.ReplaceAll(text, "\\n", "\n")
            text = regexp.MustCompile(`\{[^}]*\}`).ReplaceAllString(text, "")

						cue.Content = []subtitle.CueContentSegment{{Text: text}}
					}
				}

				sub.Cues = append(sub.Cues, cue)
			}
		case StylesSection:
			//TODO Extraction of styles
		}
	}

	return &sub, nil
}

package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SectionType int

const (
	InfoSection = iota
	StylesSection
	EventsSection
)

type LineType int

const (
	SynchPointLine = iota
	TimerLine
	DialogueLine
	FormatLine
)

func parseLineType(line string) (LineType, string, error) {
	if strings.HasPrefix(line, "Synch Point: ") {
		return SynchPointLine, strings.TrimPrefix(line, "Synch Point: "), nil
	} else if strings.HasPrefix(line, "Timer: ") {
		return TimerLine, strings.TrimPrefix(line, "Timer: "), nil
	} else if strings.HasPrefix(line, "Dialogue: ") {
		return DialogueLine, strings.TrimPrefix(line, "Dialogue: "), nil
	} else if strings.HasPrefix(line, "Format: ") {
		return FormatLine, strings.TrimPrefix(line, "Format: "), nil
	}

	return DialogueLine, "", fmt.Errorf("Invalid line type: %s", line)
}

type FieldType int

const (
	NameField = iota
	FontnameField
	FontsizeField
	PrimaryColourField
	SecondaryColourField
	TertiaryColourField
	BackColourField
	BoldField
	ItalicField
	UnderlineField
	StrikeOutField
	ScaleXField
	ScaleYField
	SpacingField
	AngleField
	BorderStyleField
	OutlineField
	ShadowField
	AlignmentField
	MarginLField
	MarginRField
	MarginVField
	AlphaLevelField
	EncodingField
	MarkedField
	StartField
	EndField
	StyleField
	EffectField
	TextField
	Unknown
)

func parseFieldType(str string) FieldType {
	switch strings.TrimSpace(str) {
	case "Name":
		return NameField
	case "Fontname":
		return FontnameField
	case "Fontsize":
		return FontsizeField
	case "PrimaryColour":
		return PrimaryColourField
	case "SecondaryColour":
		return SecondaryColourField
	case "TertiaryColour":
		return TertiaryColourField
	case "BackColour":
		return BackColourField
	case "Bold":
		return BoldField
	case "Italic":
		return ItalicField
	case "Underline":
		return UnderlineField
	case "StrikeOut":
		return StrikeOutField
	case "ScaleX":
		return ScaleXField
	case "ScaleY":
		return ScaleYField
	case "Spacing":
		return SpacingField
	case "Angle":
		return AngleField
	case "BorderStyle":
		return BorderStyleField
	case "Outline":
		return OutlineField
	case "Shadow":
		return ShadowField
	case "Alignment":
		return AlignmentField
	case "MarginL":
		return MarginLField
	case "MarginR":
		return MarginRField
	case "MarginV":
		return MarginVField
	case "AlphaLevel":
		return AlphaLevelField
	case "Encoding":
		return EncodingField
	case "Marked":
		return MarkedField
	case "Start":
		return StartField
	case "End":
		return EndField
	case "Style":
		return StyleField
	case "Effect":
		return EffectField
	case "Text":
		return TextField
	}

	return Unknown
}

type LineReader struct {
	lines    []string
	position int
}

func LineReaderFromBytes(data []byte) *LineReader {
	str := string(data)

	lines := strings.Split(str, "\r\n")

	if len(lines) == 1 {
		lines = strings.Split(str, "\n")
	}

	return &LineReader{
		lines:    lines,
		position: 0,
	}
}

func (r *LineReader) Advance() (string, error) {
	if r.AtEnd() {
		return "", fmt.Errorf("Already reached the end")
	}

	r.position++

	return r.lines[r.position-1], nil
}

func (r *LineReader) AtEnd() bool {
	return r.position >= len(r.lines)
}

type HeaderConfig struct {
	timer      float64
	synchPoint time.Duration
}

func parseTimestamp(str string) (time.Duration, error) {
	parts := strings.Split(str, ":")

	if len(parts) < 3 {
		return time.Duration(0), fmt.Errorf("Invalid timestamp: %v", str)
	}

	hr, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid hr: %v", parts[0])
	}

	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid min: %v", parts[1])
	}

	sPt := strings.Split(parts[2], ".")

	s, err := strconv.Atoi(sPt[0])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid s: %v", parts[0])
	}

	hs, err := strconv.Atoi(sPt[1])
	if err != nil {
		return time.Duration(0), fmt.Errorf("Invalid hundreth second: %v", parts[1])
	}

	return time.Duration(hr*3600_000_000_000 + min*60_000_000_000 + s*1000_000_000 + hs*10_000_000), nil
}

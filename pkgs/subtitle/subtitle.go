package subtitle

import (
	"fmt"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/subtitle/ass"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/pgs"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/srt"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
	"github.com/nandesh-dev/subtle/pkgs/warning"
)

type Format int

const (
	PGS Format = iota
	ASS
	SRT
)

func ParseFormat(f string) (Format, error) {
	switch f {
	case "ass":
		return ASS, nil
	case "pgs":
		return PGS, nil
	case "srt":
		return SRT, nil
	}

	return ASS, fmt.Errorf("Invalid format: %v", f)
}

type Segment struct {
	start time.Duration
	end   time.Duration
	text  string
}

func NewSegment() *Segment {
	return &Segment{}
}

func (s *Segment) Start() time.Duration {
	return s.start
}

func (s *Segment) End() time.Duration {
	return s.end
}

func (s *Segment) Text() string {
	return s.text
}

func (s *Segment) SetStart(start time.Duration) {
	s.start = start
}

func (s *Segment) SetEnd(end time.Duration) {
	s.end = end
}

func (s *Segment) SetText(text string) {
	s.text = text
}

type Subtitle struct {
	segments []Segment
}

func NewSubtitle() *Subtitle {
	return &Subtitle{
		segments: make([]Segment, 0),
	}
}

func (s *Subtitle) Segments() []Segment {
	return s.segments
}

func (s *Subtitle) AddSegment(segment Segment) {
	s.segments = append(s.segments, segment)
}

func FromRawStream(rawStream RawStream) (*Subtitle, warning.WarningList, error) {
	warnings := warning.NewWarningList()
	subtitle := NewSubtitle()

	switch rawStream.Format() {
	case ASS:
		assSubtitle, wrn, err := ass.DecodeSubtitle(rawStream.Filepath(), rawStream.Index())
		warnings.Append(*wrn)
		if err != nil {
			return nil, *warnings, fmt.Errorf("Error decoding ass subtitle: %v", err)
		}

		for _, assSegment := range assSubtitle.Segments() {
			segment := NewSegment()
			segment.SetStart(assSegment.Start())
			segment.SetEnd(assSegment.End())
			segment.SetText(assSegment.Text())

			subtitle.AddSegment(*segment)
		}

	case PGS:
		pgsSubtitle, wrn, err := pgs.DecodeSubtitle(rawStream.Filepath(), rawStream.Index())
		warnings.Append(*wrn)
		if err != nil {
			return nil, *warnings, fmt.Errorf("Error decoding pgs subtitle: %v", err)
		}

		tes := tesseract.NewClient()
		defer tes.Close()

		for i, pgsSegment := range pgsSubtitle.Segments() {
			segment := NewSegment()
			segment.SetStart(pgsSegment.Start())

			if i+1 < len(pgsSubtitle.Segments()) {
				segment.SetEnd(pgsSubtitle.Segments()[i+1].Start())
			} else {
				segment.SetEnd(pgsSegment.Start() + time.Second*10)
			}
			txt := ""

			images, err := pgsSegment.Images()
			if err != nil {
				return nil, *warnings, fmt.Errorf("Error getting images from pgs subtitle: %v", err)
			}

			for _, img := range images {
				line, err := tes.ExtractTextFromImage(img, rawStream.Language())
				if err != nil {
					return nil, *warnings, fmt.Errorf("Error extracting text from image: %v", err)
				}

				txt += line
			}

			segment.SetText(txt)

			subtitle.AddSegment(*segment)
		}

	}

	return subtitle, *warnings, nil
}

type EncodedSubtitleFile struct {
	extension string
	content   []byte
}

func NewEncodedSubtitleFile(content []byte, extension string) *EncodedSubtitleFile {
	return &EncodedSubtitleFile{
		extension: extension,
		content:   content,
	}
}

func (e *EncodedSubtitleFile) Extension() string {
	return e.extension
}

func (e *EncodedSubtitleFile) Content() []byte {
	return e.content
}

type EncodedSubtitle struct {
	files []EncodedSubtitleFile
}

func NewEncodedSubtitle() *EncodedSubtitle {
	return &EncodedSubtitle{
		files: make([]EncodedSubtitleFile, 0),
	}
}

func (e *EncodedSubtitle) AddFile(file EncodedSubtitleFile) {
	e.files = append(e.files, file)
}

func (e *EncodedSubtitle) Files() []EncodedSubtitleFile {
	return e.files
}

func (s *Subtitle) Encode(format Format) (*EncodedSubtitle, error) {
	switch format {
	case SRT:
		srtStream := srt.NewStream()

		for _, segment := range s.Segments() {
			srtSegment := srt.NewSegment()

			srtSegment.SetStart(segment.Start())
			srtSegment.SetEnd(segment.End())
			srtSegment.SetText(segment.Text())

			srtStream.AddSegment(*srtSegment)
		}

		str := srt.EncodeSubtitle(*srtStream)

		srtFile := NewEncodedSubtitleFile([]byte(str), ".srt")
		encodedSubtitle := NewEncodedSubtitle()
		encodedSubtitle.AddFile(*srtFile)

		return encodedSubtitle, nil
	}

	return nil, fmt.Errorf("Unsupported output format: %v", format)
}

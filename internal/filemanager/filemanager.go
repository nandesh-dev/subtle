package filemanager

import (
	"os"
	"path/filepath"

	"github.com/nandesh-dev/subtle/internal/video"
)

type SubtitleFile struct {
	Path   string
	Format SubtitleFileFormat
}

type Directory struct {
	Path      string
	Childrens []Directory
	Videos    []video.VideoFile
	Subtitles []SubtitleFile
}

var videoFormatLookup = map[string]video.VideoFileFileFormat{
	".mp4": video.MP4,
	".mkv": video.MKV,
	".avi": video.AVI,
	".mov": video.MOV,
}

type SubtitleFileFormat int

const (
	SRT SubtitleFileFormat = iota
	ASS
	SSA
	IDX
	SUB
	PGS
)

var subtitleFormatLookup = map[string]SubtitleFileFormat{
	".srt": SRT,
	".ass": ASS,
	".ssa": SSA,
	".idx": IDX,
	".sub": SUB,
	".PGS": PGS,
}

func ReadDirectory(path string) (*Directory, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	childrenDirectories := make([]Directory, 0)
	videoFiles := make([]video.VideoFile, 0)
	subtitleFiles := make([]SubtitleFile, 0)

	for _, entry := range files {
		if entry.IsDir() {
			childrenDirectory, err := ReadDirectory(filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}

			childrenDirectories = append(childrenDirectories, *childrenDirectory)
		}

		extension := filepath.Ext(entry.Name())

		videoFormat, isVideoFile := videoFormatLookup[extension]
		if isVideoFile {
			videoFile := video.VideoFile{
				Path:   filepath.Join(path, entry.Name()),
				Format: videoFormat,
			}
			videoFiles = append(videoFiles, videoFile)
			continue
		}

		subtitleFormat, isSubtitleFile := subtitleFormatLookup[extension]
		if isSubtitleFile {
			subtitleFile := SubtitleFile{
				Path:   filepath.Join(path, entry.Name()),
				Format: subtitleFormat,
			}
			subtitleFiles = append(subtitleFiles, subtitleFile)
			continue
		}
	}

	return &Directory{
		Path:      path,
		Childrens: childrenDirectories,
		Videos:    videoFiles,
		Subtitles: subtitleFiles,
	}, nil
}

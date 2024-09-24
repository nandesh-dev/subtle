package filemanager

import (
	"os"
	"path/filepath"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/video"
)

type Directory struct {
	Path      string
	Childrens []Directory
	Videos    []video.VideoFile
	Subtitles []subtitle.SubtitleFile
}

var videoFormatLookup = map[string]video.VideoFileFileFormat{
	".mp4": video.MP4,
	".mkv": video.MKV,
	".avi": video.AVI,
	".mov": video.MOV,
}

var subtitleFormatLookup = map[string]subtitle.SubtitleFileFormat{
	".srt": subtitle.SRT,
	".ass": subtitle.ASS,
	".ssa": subtitle.SSA,
	".idx": subtitle.IDX,
	".sub": subtitle.SUB,
	".PGS": subtitle.PGS,
}

func ReadDirectory(path string) (*Directory, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	childrenDirectories := make([]Directory, 0)
	videoFiles := make([]video.VideoFile, 0)
	subtitleFiles := make([]subtitle.SubtitleFile, 0)

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
			subtitleFile := subtitle.SubtitleFile{
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
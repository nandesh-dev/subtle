package extract

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"slices"
	"strings"

	"github.com/nandesh-dev/subtle/pkgs/ass"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/pgs"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
	"github.com/nandesh-dev/subtle/pkgs/warning"
	"golang.org/x/text/language"
)

func Run() warning.WarningList {
	warnings := warning.NewWarningList()
	for _, rootDirectoryConfig := range config.Config().Media.RootDirectories {
		dir, _, err := filemanager.ReadDirectory(rootDirectoryConfig.Path)
		if err != nil {
			warnings.AddWarning(fmt.Errorf("Error reading root directory: %v; %v", rootDirectoryConfig.Path, err))
			continue
		}

		warns := extractSubtitleFromDirectory(*dir, rootDirectoryConfig.AutoExtract)
		warnings.Append(warns)
	}

	return *warnings
}

func extractSubtitleFromDirectory(dir filemanager.Directory, autoExtractConfig config.AutoExtract) warning.WarningList {
	warnings := warning.NewWarningList()

	for _, video := range dir.VideoFiles() {
		var videoEntry db.Video
		if err := db.DB().Where(&db.Video{DirectoryPath: video.DirectoryPath(), Filename: video.Filename()}).
			Preload("Subtitles").
			Preload("Subtitles.Segments").
			First(&videoEntry).Error; err != nil {
			log.Fatal("Error getting entry: ", err, video.DirectoryPath(), video.Filename())
		}

		rawStreams, err := video.RawStreams()
		if err != nil {
			warnings.AddWarning(fmt.Errorf("Error extracting raw stream from video: %v; %v", video.Filepath(), err))
			continue
		}

		rawStreamRanks := map[filemanager.RawStream]int{}

		for _, rawStream := range *rawStreams {
			if !slices.ContainsFunc(videoEntry.Subtitles, func(subtitleEntry db.Subtitle) bool {
				subtitleEntryLanguageTag, _ := language.Parse(subtitleEntry.Language)
				return subtitleEntryLanguageTag == rawStream.Language()
			}) && slices.ContainsFunc(autoExtractConfig.Languages, func(lang language.Tag) bool {
				return lang == rawStream.Language()
			}) {
				rawStreamRanks[rawStream] = 1

				for _, titleKeyword := range autoExtractConfig.RawStreamTitleKeywords {
					if strings.Contains(rawStream.Title(), titleKeyword) {
						rawStreamRanks[rawStream]++
					}
				}
			}
		}

		highestRank := 0
		var highestRankRawStream filemanager.RawStream
		for rawStream, rank := range rawStreamRanks {
			if rank >= highestRank {
				highestRank = rank
				highestRankRawStream = rawStream
			}
		}

		if highestRank == 0 {
			continue
		}

		rawStream := highestRankRawStream
		var sub subtitle.Subtitle

		switch rawStream.Format() {
		case subtitle.ASS:
			s, warns, err := ass.DecodeSubtitle(rawStream)
			warnings.Append(warns)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error decoding subtitle for video: %v; %v", video.Filepath(), err))
				continue
			}

			sub = *s

		case subtitle.PGS:
			s, warns, err := pgs.DecodeSubtitle(rawStream)
			warnings.Append(warns)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error decoding subtitle for video: %v; %v", video.Filepath(), err))
				continue
			}

			sub = *s
		}

		if sub == nil {
			continue
		}

		switch sub := sub.(type) {
		case subtitle.TextSubtitle:
			subtitleEntry := db.Subtitle{
				Language: rawStream.Language().String(),
				Segments: make([]db.Segment, 0),
			}

			for _, segment := range sub.Segments() {
				segmentEntry := db.Segment{
					StartTime:    segment.Start(),
					EndTime:      segment.End(),
					Text:         segment.Text(),
					OriginalText: segment.Text(),
				}

				subtitleEntry.Segments = append(subtitleEntry.Segments, segmentEntry)
			}

			videoEntry.Subtitles = append(videoEntry.Subtitles, subtitleEntry)

			db.DB().Save(&videoEntry)
		case subtitle.ImageSubtitle:
			subtitleEntry := db.Subtitle{
				Language: rawStream.Language().String(),
				Segments: make([]db.Segment, 0),
			}

			tesseractClient := tesseract.NewClient()
			defer tesseractClient.Close()

			for _, segment := range sub.Segments() {
				imageDataBuffer := new(bytes.Buffer)
				if err := png.Encode(imageDataBuffer, segment.Image()); err != nil {
					warnings.AddWarning(fmt.Errorf("Error encoding image to png for video: %v; %v", video.Filepath(), err))
					continue
				}

				text, err := tesseractClient.ExtractTextFromPNGImage(*imageDataBuffer, rawStream.Language())
				if err != nil {
					warnings.AddWarning(fmt.Errorf("Error extracting text from image: %v", err))
				}

				segmentEntry := db.Segment{
					StartTime:     segment.Start(),
					EndTime:       segment.End(),
					Text:          text,
					OriginalText:  text,
					OriginalImage: imageDataBuffer.Bytes(),
				}

				subtitleEntry.Segments = append(subtitleEntry.Segments, segmentEntry)
			}

			videoEntry.Subtitles = append(videoEntry.Subtitles, subtitleEntry)

			db.DB().Save(&videoEntry)
		}
	}

	return *warnings
}

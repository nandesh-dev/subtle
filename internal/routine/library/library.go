package library

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/warning"
	"golang.org/x/text/language"
)

func RunLibraryRoutine() {
	warnings := warning.NewWarningList()

	for _, watchDirectoryConfig := range config.Config().Media.WatchDirectories {
		dir, _, err := filemanager.ReadDirectory(watchDirectoryConfig.Path)
		if err != nil {
			return
		}

		for _, videoFile := range dir.VideoFiles() {
			missingSubtitleLanguages := make([]language.Tag, 0)
			for _, languageCode := range watchDirectoryConfig.AutoExtract.Languages {
				languageTag, err := language.Parse(languageCode)
				if err != nil {
					warnings.AddWarning(fmt.Errorf("Invalid language code in config: %v; %v", languageCode, err))
					continue
				}

				hasSubtitleLanguage, wrn := videoFile.HasSubtitleLanguage(languageTag)
				warnings.Append(wrn)

				if !hasSubtitleLanguage {
					missingSubtitleLanguages = append(missingSubtitleLanguages, languageTag)
				}
			}

			if len(missingSubtitleLanguages) == 0 {
				continue
			}

			rawStreams, err := subtitle.ExtractRawStreams(videoFile.Filepath())
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error extracting raw streams from video file: %v", err))
			}

			for _, formatCode := range watchDirectoryConfig.AutoExtract.Formats {
				format, err := subtitle.ParseFormat(formatCode)
				if err != nil {
					warnings.AddWarning(fmt.Errorf("Invalid subtitle format in config: %v", formatCode))
					continue
				}

				for _, rawStream := range rawStreams {
					if format == rawStream.Format() && slices.Contains(missingSubtitleLanguages, rawStream.Language()) {
						sub, wrn, err := subtitle.FromRawStream(rawStream)
						warnings.Append(wrn)
						if err != nil {
							warnings.AddWarning(fmt.Errorf("Error extracting subtitle: %v", err))
							return
						}

						outputFormat, err := subtitle.ParseFormat(watchDirectoryConfig.AutoExtract.OutputFormat)
						if err != nil {
							warnings.AddWarning(fmt.Errorf("Invalid output format in config: %v", err))
							return
						}

						encodedSubtitle, _ := sub.Encode(outputFormat)

						for _, file := range encodedSubtitle.Files() {
							path := filepath.Join(videoFile.DirectoryPath(), videoFile.Basename()+"."+rawStream.Language().String()) + file.Extension()

							if err := os.WriteFile(path, file.Content(), 0644); err != nil {
								fmt.Print(err)
							}
						}
					}
				}
			}
		}
	}

	fmt.Println("WARNINGS BEGIN")
	for _, warning := range warnings.Warnings() {
		fmt.Println(warning)
	}
}

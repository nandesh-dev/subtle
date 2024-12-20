package extract

import (
	"bytes"
	"context"
	"image/png"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/ass"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/pgs"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language"
)

func Run(conf *config.Config, db *ent.Client) {
	logger := logging.NewRoutineLogger("extract")

	c, err := conf.Read()
	if err != nil {
		logger.Error("cannot read config", "err", err)
		return
	}

	for _, mediaDirectoryConfig := range c.MediaDirectories {
		if !mediaDirectoryConfig.Extraction.Enable {
			continue
		}

		videoEntries, err := db.Video.Query().Where(video_schema.FilepathHasPrefix(mediaDirectoryConfig.Path)).All(context.Background())
		if err != nil {
			logger.Error("cannot get videos from database", "err", err)
			continue
		}

		for _, videoEntry := range videoEntries {
			logger := logger.With("video_filepath", videoEntry.Filepath)

			// Skip if a subtitle is already extracted for the video
			extractedSubtitleCount, err := videoEntry.QuerySubtitles().Where(subtitle_schema.Extracted(true)).Count(context.Background())
			if err != nil {
				logger.Error("cannot count extracted subtitle in database", "err", err)
				continue
			} else {
				if extractedSubtitleCount > 0 {
					continue
				}
			}

			logger.Info("looking for suitable subtitle to extract")

			// Get subtitles from database which contain the required language for the format
			searchImportFormatStrings := []string{}
			if mediaDirectoryConfig.Extraction.Formats.ASS.Enable {
				searchImportFormatStrings = append(searchImportFormatStrings, subtitle.MapFormat(subtitle.ASS))
			}
			if mediaDirectoryConfig.Extraction.Formats.PGS.Enable {
				searchImportFormatStrings = append(searchImportFormatStrings, subtitle.MapFormat(subtitle.PGS))
			}

			requiredASSLanguageStrings := make([]string, 0)
			for _, lang := range mediaDirectoryConfig.Extraction.Formats.ASS.Languages {
				requiredASSLanguageStrings = append(requiredASSLanguageStrings, lang.String())
			}

			requiredPGSLanguageStrings := make([]string, 0)
			for _, lang := range mediaDirectoryConfig.Extraction.Formats.ASS.Languages {
				requiredPGSLanguageStrings = append(requiredPGSLanguageStrings, lang.String())
			}

			subtitleEntries, err := videoEntry.QuerySubtitles().Where(
				subtitle_schema.And(
					subtitle_schema.ImportFormatIn(searchImportFormatStrings...),
					subtitle_schema.Or(
						subtitle_schema.And(
							subtitle_schema.ImportFormat(subtitle.MapFormat(subtitle.ASS)),
							subtitle_schema.LanguageIn(requiredASSLanguageStrings...),
						),
						subtitle_schema.And(
							subtitle_schema.ImportFormat(subtitle.MapFormat(subtitle.PGS)),
							subtitle_schema.LanguageIn(requiredPGSLanguageStrings...),
						),
					),
				)).All(context.Background())
			if err != nil {
				logger.Error("cannot get subitles from database", "err", err)
				continue
			}

			//Find the best subtitle based on the title keywords
			bestScore := -1
			bestScoreSubtitleId := 0

			for _, subtitleEntry := range subtitleEntries {
				score := 0

				for _, rawStreamTitleKeyword := range mediaDirectoryConfig.Extraction.RawStreamTitleKeywords {
					if strings.Contains(subtitleEntry.Title, rawStreamTitleKeyword) {
						score++
					}
				}

				if score > bestScore {
					bestScore = score
					bestScoreSubtitleId = subtitleEntry.ID
				}
			}

			if bestScore == -1 {
				logger.Info("no suitable subtitle found")
				continue
			}

			subtitleEntry, err := db.Subtitle.Query().Where(subtitle_schema.ID(bestScoreSubtitleId)).Only(context.Background())
			if err != nil {
				logger.Error("cannot get subtitle from database", "err", err)
				continue
			}

			logger = logger.With("subtitle_title", subtitleEntry.Title)

			// Skip if the subtitle is currently being processing
			if subtitleEntry.Processing {
				logger.Info("skipping! subtitle is already being processed")
				continue
			}

			logger.Info("subtitle found, processing it")

			if err := db.Subtitle.UpdateOne(subtitleEntry).SetProcessing(true).Exec(context.Background()); err != nil {
				logger.Error("cannot update subtitle processing status", "err", err)
				continue
			}

			defer func() {
				if err := db.Subtitle.UpdateOneID(subtitleEntry.ID).SetProcessing(false).Exec(context.Background()); err != nil {
					logger.Error("cannot update subtitle processing status", "err", err)
				}
			}()

			format, err := subtitle.ParseFormat(subtitleEntry.ImportFormat)
			if err != nil {
				logger.Error("cannot parse subtitle format", "err", err)
				continue
			}

			lang, err := language.Parse(subtitleEntry.Language)
			if err != nil {
				logger.Error("cannot parse subtitle language", "err", err)
				continue
			}

			rawStream := filemanager.NewRawStream(
				videoEntry.Filepath,
				int(subtitleEntry.ImportVideoStreamIndex),
				format,
				lang,
				subtitleEntry.Title,
			)

			var sub subtitle.Subtitle

			logger.Info("extracting subtitle from raw stream")

			logger = logger.With("subtitle_format", subtitleEntry.ImportFormat)

			switch format {
			case subtitle.ASS:
				s, _, err := ass.ExtractFromRawStream(*rawStream)
				if err != nil {
					logger.Error("cannot extract subtitle", "err", err)
					continue
				}

				sub = *s
			case subtitle.PGS:
				s, _, err := pgs.ExtractFromRawStream(*rawStream)
				if err != nil {
					logger.Error("cannot extract subtitle", "err", err)
					continue
				}

				sub = *s

			default:
				logger.Error("unsupported / invalid subtitle format")
				continue
			}

			tx, err := db.Tx(context.Background())
			if err != nil {
				logger.Error(logging.DatabaseTransactionCreateError, "err", err)
				continue
			}

			if err := func() error {
				switch sub := sub.(type) {
				case subtitle.TextSubtitle:
					for _, segment := range sub.Segments() {
						if err := tx.Segment.Create().
							AddSubtitle(subtitleEntry).
							SetStartTime(segment.Start()).
							SetEndTime(segment.End()).
							SetOriginalText(segment.Text()).
							Exec(context.Background()); err != nil {
							logger.Error("cannot add segment to subtitle", "err", err)
							return err
						}
					}
				case subtitle.ImageSubtitle:
					for _, segment := range sub.Segments() {
						imageDataBuffer := new(bytes.Buffer)
						if err := png.Encode(imageDataBuffer, segment.Image()); err != nil {
							logger.Error("cannot encode image to png", "err", err)
							return err
						}

						if err := tx.Segment.Create().
							AddSubtitle(subtitleEntry).
							SetStartTime(segment.Start()).
							SetEndTime(segment.End()).
							SetOriginalImage(imageDataBuffer.Bytes()).
							Exec(context.Background()); err != nil {
							logger.Error("cannot add segment to subtitle", "err", err)
							return err
						}
					}
				}

				if err := tx.Subtitle.UpdateOne(subtitleEntry).SetExtracted(true).Exec(context.Background()); err != nil {
					logger.Error("cannot update subtitle to extracted", "err", err)
					return err
				}

				logger.Info("subtitle extracted")

				return nil
			}(); err != nil {
				if err := tx.Rollback(); err != nil {
					logger.Error(logging.DatabaseTransactionRollbackError, "err", err)
				}
			} else {
				if err := tx.Commit(); err != nil {
					logger.Error(logging.DatabaseTransactionCommitError, "err", err)
				}
			}
		}
	}
}

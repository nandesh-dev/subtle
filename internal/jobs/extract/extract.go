package extract

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log/slog"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/ass"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/pgs"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/text/language"
)

func Run(logger *slog.Logger, conf *config.Config, db *ent.Client) {
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
			extractedSubtitleCount, err := videoEntry.QuerySubtitles().
				Where(subtitle_schema.StageNEQ(subtitle_schema.StageDetected)).
				Count(context.Background())
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

			if subtitleEntry.IsProcessing {
				logger.Warn("skipping! subtitle is already being processed")
				continue
			}

			logger.Info("subtitle found, processing it")

			if err := db.Subtitle.UpdateOne(subtitleEntry).SetIsProcessing(true).Exec(context.Background()); err != nil {
				logger.Error("cannot update subtitle processing status", "err", err)
				continue
			}

			defer func() {
				if err := db.Subtitle.UpdateOneID(subtitleEntry.ID).SetIsProcessing(false).Exec(context.Background()); err != nil {
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

			logger.Info("extracting subtitle from raw stream")
			logger = logger.With("subtitle_format", subtitleEntry.ImportFormat)

			formatString := "ass"
			if format == subtitle.PGS {
				formatString = "sup"
			}

			var subtitleBuffer, errorBuffer bytes.Buffer
			ffmpeg.LogCompiledCommand = false
			if err := ffmpeg.Input(rawStream.Filepath()).
				Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", rawStream.Index()), "c:s": "copy", "f": formatString}).
				WithOutput(&subtitleBuffer).
				WithErrorOutput(&errorBuffer).
				Run(); err != nil {
				logger.Error("cannot extract subtitle from video", "err", err, "errorBuffer", errorBuffer.String())
				continue
			}

			var parser subtitle.Parser
			switch format {
			case subtitle.ASS:
				parser = ass.NewParser()
			case subtitle.PGS:
				parser = pgs.NewParser()
			default:
				logger.Error("unsupported / invalid subtitle format")
				continue
			}

			sub, err := parser.Parse(subtitleBuffer.Bytes())
			if err != nil {
				logger.Error("cannot parse subtitle", slog.String("format", subtitleEntry.ImportFormat), slog.Any("err", err))
				continue
			}

			tx, err := db.Tx(context.Background())
			if err != nil {
				logger.Error(logging.DatabaseTransactionCreateError, "err", err)
				continue
			}

			if err := func() error {
				for _, cue := range sub.Cues {
					cueEntry := tx.Cue.Create().
						AddSubtitle(subtitleEntry).
						SetTimestampStart(cue.Timestamp.Start).
						SetTimestampEnd(cue.Timestamp.End)

					for i, originalImage := range cue.OriginalImages {
						var imageDataBuffer bytes.Buffer
						if err := png.Encode(&imageDataBuffer, originalImage); err != nil {
							logger.Error("cannot encode original image to png", "err", err)
							return err
						}

						originalImageEntry, err := tx.CueOriginalImage.Create().
							SetPosition(int32(i)).
							SetData(imageDataBuffer.Bytes()).
							Save(context.Background())
						if err != nil {
							logger.Error("cannot save original image to database", "err", err)
							return err
						}

						cueEntry.AddCueOriginalImages(originalImageEntry)
					}

					for i, content := range cue.Content {
						contentSegment, err := tx.CueContentSegment.Create().
							SetPosition(int32(i)).
							SetText(content.Text).
							Save(context.Background())
						if err != nil {
							logger.Error("cannot save content segment to database", "err", err)
							return err
						}

						cueEntry.AddCueContentSegments(contentSegment)
					}

					if err := cueEntry.
						Exec(context.Background()); err != nil {
						logger.Error("cannot add cue to subtitle", "err", err)
						return err
					}
				}

				if err := tx.Subtitle.UpdateOne(subtitleEntry).SetStage(subtitle_schema.StageExtracted).Exec(context.Background()); err != nil {
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

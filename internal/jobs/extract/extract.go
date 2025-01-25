package extract

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log/slog"
	"sort"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/ass"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/pgs"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Run(ctx context.Context, logger *slog.Logger, configFile *configuration.File, db *ent.Client) {
	config, err := configFile.Read()
	if err != nil {
		logger.Error("cannot read config", "err", err)
		return
	}

	offset := 0
	for {
		videoEntries, err := db.VideoSchema.Query().Limit(10).Offset(offset).All(ctx)
		if err != nil {
			logger.Error("cannot get videos from database", "err", err)
			return
		}

		if len(videoEntries) == 0 {
			break
		}

		offset += len(videoEntries)

		for _, videoEntry := range videoEntries {
			handleVideo(videoEntry, ctx, logger, config, db)
		}
	}
}

func handleVideo(videoEntry *ent.VideoSchema, ctx context.Context, logger *slog.Logger, config *configuration.Config, db *ent.Client) {
	logger = logger.With("video_filepath", videoEntry.Filepath)

	logger.Info("checking subtitles")
	for i, groupConfig := range config.Job.Extracting {
		logger := logger.With("config_group_index", i)

		limitInUse := groupConfig.Limit > 0

		logger.Info("counting extracted subtitles meeting group condition")
		extractedSubtitleCountQuery := videoEntry.QuerySubtitles().
			Where(subtitleschema.StageNEQ(subtitleschema.StageDetected))

		if len(groupConfig.Condition.Formats) > 0 {
			extractedSubtitleCountQuery = extractedSubtitleCountQuery.Where(subtitleschema.ImportFormatIn(groupConfig.Condition.Formats...))
		}

		if len(groupConfig.Condition.Languages) > 0 {
			extractedSubtitleCountQuery = extractedSubtitleCountQuery.Where(subtitleschema.LanguageIn(groupConfig.Condition.Languages...))
		}

		extractedSubtitleEntryCount, err := extractedSubtitleCountQuery.Count(ctx)
		if err != nil {
			logger.Error("cannot get extracted subtitle count from database", "err", err)
			return
		}

		if limitInUse && extractedSubtitleEntryCount >= groupConfig.Limit {
			logger.Info("enough subtitles are already extracted for the group; skipping")
			continue
		}

		logger.Info("looking for more subtitles to extract")
		detectedSubtitlesQuery := videoEntry.QuerySubtitles().
			Where(subtitleschema.StageEQ(subtitleschema.StageDetected))

		if len(groupConfig.Condition.Formats) > 0 {
			detectedSubtitlesQuery = detectedSubtitlesQuery.Where(subtitleschema.ImportFormatIn(groupConfig.Condition.Formats...))
		}

		if len(groupConfig.Condition.Languages) > 0 {
			detectedSubtitlesQuery = detectedSubtitlesQuery.Where(subtitleschema.LanguageIn(groupConfig.Condition.Languages...))
		}

		detectedSubtitleEntries, err := detectedSubtitlesQuery.All(ctx)
		if err != nil {
			logger.Error("cannot get subtitles from database", "err", err)
		}

		if limitInUse && len(detectedSubtitleEntries) > groupConfig.Limit-extractedSubtitleEntryCount {
			logger.Info("sorting subtitles according to priorities")
			sort.Slice(detectedSubtitleEntries, func(a, b int) bool {
				return scoreSubtitle(detectedSubtitleEntries[a], &groupConfig) > scoreSubtitle(detectedSubtitleEntries[b], &groupConfig)
			})
		}

		for _, detectedSubtitleEntry := range detectedSubtitleEntries {
			if limitInUse && extractedSubtitleEntryCount >= groupConfig.Limit {
				break
			}

			handleSubtitle(detectedSubtitleEntry, ctx, logger, db)
			extractedSubtitleEntryCount++
		}
	}
}

func scoreSubtitle(subtitleEntry *ent.SubtitleSchema, groupConfig *configuration.ExtractingGroup) int {
	totalScore := 0

	if score, exist := groupConfig.Priority.Format[subtitleEntry.ImportFormat]; exist {
		totalScore += score
	}

	if score, exist := groupConfig.Priority.Language[subtitleEntry.Language]; exist {
		totalScore += score
	}

	for nameKeyword, score := range groupConfig.Priority.TitleKeyword {
		if strings.Contains(subtitleEntry.Title, nameKeyword) {
			totalScore += score
		}
	}

	return totalScore
}

func handleSubtitle(subtitleEntry *ent.SubtitleSchema, ctx context.Context, logger *slog.Logger, db *ent.Client) {
	logger.With("subtitle_title", subtitleEntry.Title)

	if subtitleEntry.IsProcessing {
		logger.Warn("subtitle is already being processed; skipping")
		return
	}

	logger.Info("extracting subtitle from raw stream")

	if err := subtitleEntry.Update().SetIsProcessing(true).Exec(ctx); err != nil {
		logger.Error("cannot update subtitle processing status", "err", err)
		return
	}

	defer func() {
		if err := subtitleEntry.Update().SetIsProcessing(false).Exec(ctx); err != nil {
			logger.Error("cannot update subtitle processing status", "err", err)
		}
	}()

	videoEntry, err := subtitleEntry.QueryVideo().Only(ctx)
	if err != nil {
		logger.Error("cannot get subtitle video from database", "err", err)
		return
	}

	var subtitleBuffer, errorBuffer bytes.Buffer
	ffmpeg.LogCompiledCommand = false
	if err := ffmpeg.Input(videoEntry.Filepath).
		Output("pipe:", ffmpeg.KwArgs{
			"map": fmt.Sprintf("0:%v", *subtitleEntry.ImportVideoStreamIndex),
			"c:s": "copy",
			"f":   subtitleEntry.ImportFormat.FFMpegString(),
		}).
		WithOutput(&subtitleBuffer).
		WithErrorOutput(&errorBuffer).
		Run(); err != nil {
		logger.Error("cannot extract subtitle from video", "err", err, "errorBuffer", errorBuffer.String())
		return
	}

	var parser subtitle.Parser
	switch subtitleEntry.ImportFormat {
	case subtitle.ASS:
		parser = ass.NewParser()
	case subtitle.PGS:
		parser = pgs.NewParser()
	default:
		logger.Error("unsupported / invalid subtitle format")
		return
	}

	logger = logger.With("format", subtitleEntry.ImportFormat.String())

	sub, err := parser.Parse(subtitleBuffer.Bytes())
	if err != nil {
		logger.Error("cannot parse subtitle", "err", err)
		return
	}

	tx, err := db.Tx(ctx)
	if err != nil {
		logger.Error(logging.DatabaseTransactionCreateError, "err", err)
		return
	}

	if err := func() error {
		for _, cue := range sub.Cues {
			cueEntryQuery := tx.SubtitleCueSchema.Create().
				AddSubtitle(subtitleEntry).
				SetTimestampStart(cue.Timestamp.Start).
				SetTimestampEnd(cue.Timestamp.End)

			for i, originalImage := range cue.OriginalImages {
				var imageDataBuffer bytes.Buffer
				if err := png.Encode(&imageDataBuffer, originalImage); err != nil {
					logger.Error("cannot encode original image to png", "err", err)
					return err
				}

				originalImageEntry, err := tx.SubtitleCueOriginalImageSchema.Create().
					SetPosition(int32(i)).
					SetData(imageDataBuffer.Bytes()).
					Save(ctx)
				if err != nil {
					logger.Error("cannot save original image to database", "err", err)
					return err
				}

				cueEntryQuery.AddOriginalImages(originalImageEntry)
			}

			for i, content := range cue.Content {
				contentSegment, err := tx.SubtitleCueContentSegmentSchema.Create().
					SetPosition(i).
					SetText(content.Text).
					Save(ctx)
				if err != nil {
					logger.Error("cannot save content segment to database", "err", err)
					return err
				}

				cueEntryQuery.AddContentSegments(contentSegment)
			}

			if err := cueEntryQuery.
				Exec(ctx); err != nil {
				logger.Error("cannot add cue to subtitle", "err", err)
				return err
			}
		}

		if err := tx.SubtitleSchema.UpdateOne(subtitleEntry).
			SetStage(subtitleschema.StageExtracted).
			Exec(ctx); err != nil {
			logger.Error("cannot update subtitle to extracted", "err", err)
			return err
		}

		logger.Info("subtitle extracted")
		return nil
	}(); err != nil {
		logger.Warn("rolling back")

		if err := tx.Rollback(); err != nil {
			logger.Error(logging.DatabaseTransactionRollbackError, "err", err)
		}
	} else {
		if err := tx.Commit(); err != nil {
			logger.Error(logging.DatabaseTransactionCommitError, "err", err)
		}
	}
}

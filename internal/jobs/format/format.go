package format

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/subtitlecueschema"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
	"github.com/nandesh-dev/subtle/pkgs/logging"
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

	logger.Info("checking video")
	for i, groupConfig := range config.Job.Formating {
		logger := logger.With("config_group_index", i)

		limitInUse := groupConfig.Limit > 0

		logger.Info("counting formated subtitles meeting group condition")
		formatedSubtitleCountQuery := videoEntry.QuerySubtitles().
			Where(subtitleschema.Or(
				subtitleschema.StageEQ(subtitleschema.StageFormated),
				subtitleschema.StageEQ(subtitleschema.StageExported)))

		if len(groupConfig.Condition.Formats) > 0 {
			formatedSubtitleCountQuery = formatedSubtitleCountQuery.Where(subtitleschema.ImportFormatIn(groupConfig.Condition.Formats...))
		}

		if len(groupConfig.Condition.Languages) > 0 {
			formatedSubtitleCountQuery = formatedSubtitleCountQuery.Where(subtitleschema.LanguageIn(groupConfig.Condition.Languages...))
		}

		formatedSubtitleEntryCount, err := formatedSubtitleCountQuery.Count(ctx)
		if err != nil {
			logger.Error("cannot get formated subtitle count from database", "err", err)
			return
		}

		if limitInUse && formatedSubtitleEntryCount >= groupConfig.Limit {
			logger.Info("enough subtitles are already formated for the group; skipping")
			continue
		}

		logger.Info("looking for more subtitles to format")
		extractedSubtitlesQuery := videoEntry.QuerySubtitles().
			Where(subtitleschema.StageEQ(subtitleschema.StageExtracted))

		if len(groupConfig.Condition.Formats) > 0 {
			extractedSubtitlesQuery = extractedSubtitlesQuery.Where(subtitleschema.ImportFormatIn(groupConfig.Condition.Formats...))
		}
		if len(groupConfig.Condition.Languages) > 0 {
			extractedSubtitlesQuery = extractedSubtitlesQuery.Where(subtitleschema.LanguageIn(groupConfig.Condition.Languages...))
		}

		extractedSubtitleEntries, err := extractedSubtitlesQuery.All(ctx)
		if err != nil {
			logger.Error("cannot get subtitles from database", "err", err)
		}

		logger.Info("sorting subtitles according to priorities")
		sort.Slice(extractedSubtitleEntries, func(a, b int) bool {
			return scoreSubtitle(extractedSubtitleEntries[a], &groupConfig) > scoreSubtitle(extractedSubtitleEntries[b], &groupConfig)
		})

		for _, extractedSubtitleEntry := range extractedSubtitleEntries {
			if limitInUse && formatedSubtitleEntryCount >= groupConfig.Limit {
				break
			}

			handleSubtitle(extractedSubtitleEntry, ctx, logger, db, &groupConfig)
			formatedSubtitleEntryCount++
		}
	}
}

func scoreSubtitle(subtitleEntry *ent.SubtitleSchema, groupConfig *configuration.FormatingGroup) int {
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

func handleSubtitle(subtitleEntry *ent.SubtitleSchema, ctx context.Context, logger *slog.Logger, db *ent.Client, groupConfig *configuration.FormatingGroup) {
	logger = logger.With("subtitle_title", subtitleEntry.Title)

	logger.Info("formating subtitle")

	if err := subtitleEntry.Update().SetIsProcessing(true).Exec(ctx); err != nil {
		logger.Error("cannot update subtitle processing status", "err", err)
		return
	}

	defer func() {
		if err := subtitleEntry.Update().SetIsProcessing(false).Exec(ctx); err != nil {
			logger.Error("cannot update subtitle processing status", "err", err)
		}
	}()

	tx, err := db.Tx(ctx)
	if err != nil {
		logger.Error(logging.DatabaseTransactionCreateError, "err", err)
		return
	}

	cueEntries, err := tx.SubtitleCueSchema.Query().
		Where(subtitlecueschema.HasSubtitleWith(subtitleschema.ID(subtitleEntry.ID))).
		All(ctx)
	if err != nil {
		logger.Error("cannot get subtitle cues", "err", err)
		return
	}

	if err := func() error {
		for _, cueEntry := range cueEntries {
			subtitleCueContentSegmentEntries, err := cueEntry.QueryContentSegments().All(ctx)
			if err != nil {
				logger.Error("cannot get subtitle cue content segments from database", "err", err)
				return err
			}

			for _, segmentEntry := range subtitleCueContentSegmentEntries {
				text := segmentEntry.Text

				for _, mapping := range groupConfig.Config.WordMappings {
					text = strings.ReplaceAll(text, mapping.From, mapping.To)
				}

				if err := tx.SubtitleCueContentSegmentSchema.UpdateOne(segmentEntry).
					SetText(text).
					Exec(ctx); err != nil {
					logger.Error("cannot update content segment", "err", err)
					return err
				}
			}
		}

		if err := tx.SubtitleSchema.UpdateOne(subtitleEntry).
			SetStage(subtitleschema.StageFormated).
			Exec(ctx); err != nil {
			logger.Error("cannot udpate subtitle stage", "err", err)
			return err
		}

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

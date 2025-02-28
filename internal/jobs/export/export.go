package export

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/srt"
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
			handleVideo(videoEntry, ctx, logger, config)
		}
	}
}

func handleVideo(videoEntry *ent.VideoSchema, ctx context.Context, logger *slog.Logger, config *configuration.Config) {
	logger = logger.With("video_filepath", videoEntry.Filepath)

	for i, groupConfig := range config.Job.Exporting {
		logger := logger.With("config_group_index", i)

		limitInUse := groupConfig.Limit > 0

		logger.Info("counting formated subtitles meeting group conditions")
		exportedSubtitleCountQuery := videoEntry.QuerySubtitles().
			Where(subtitleschema.StageEQ(subtitleschema.StageExported))

		if len(groupConfig.Condition.Formats) > 0 {
			exportedSubtitleCountQuery = exportedSubtitleCountQuery.Where(subtitleschema.ImportFormatIn(groupConfig.Condition.Formats...))
		}

		if len(groupConfig.Condition.Languages) > 0 {
			exportedSubtitleCountQuery = exportedSubtitleCountQuery.Where(subtitleschema.LanguageIn(groupConfig.Condition.Languages...))
		}

		exportedSubtitleEntryCount, err := exportedSubtitleCountQuery.Count(ctx)
		if err != nil {
			logger.Error("cannot get formated subtitle count from database", "err", err)
			return
		}

		if limitInUse && exportedSubtitleEntryCount >= groupConfig.Limit {
			logger.Info("enough subtitles are already exported for the group; skipping")
			continue
		}

		logger.Info("looking for more subtitles to export")
		formatedSubtitleQuery := videoEntry.QuerySubtitles().
			Where(subtitleschema.StageEQ(subtitleschema.StageFormated))

		if len(groupConfig.Condition.Formats) > 0 {
			formatedSubtitleQuery = formatedSubtitleQuery.Where(subtitleschema.ImportFormatIn(groupConfig.Condition.Formats...))
		}

		if len(groupConfig.Condition.Languages) > 0 {
			formatedSubtitleQuery = formatedSubtitleQuery.Where(subtitleschema.LanguageIn(groupConfig.Condition.Languages...))
		}

		formatedSubtitleEntries, err := formatedSubtitleQuery.All(ctx)
		if err != nil {
			logger.Error("cannot get subtitles from database", "err", err)
		}

		logger.Info("sorting subtitles according to priorities")
		sort.Slice(formatedSubtitleEntries, func(a, b int) bool {
			return scoreSubtitle(formatedSubtitleEntries[a], &groupConfig) > scoreSubtitle(formatedSubtitleEntries[b], &groupConfig)
		})

		for _, formatedSubtitleEntry := range formatedSubtitleEntries {
			if limitInUse && exportedSubtitleEntryCount >= groupConfig.Limit {
				break
			}

			handleSubtitle(formatedSubtitleEntry, ctx, logger, &groupConfig)
			exportedSubtitleEntryCount++
		}
	}

}

func scoreSubtitle(subtitleEntry *ent.SubtitleSchema, groupConfig *configuration.ExportingGroup) int {
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

func handleSubtitle(subtitleEntry *ent.SubtitleSchema, ctx context.Context, logger *slog.Logger, groupConfig *configuration.ExportingGroup) {
	logger = logger.With("subtitle_title", subtitleEntry.Title)

	if groupConfig.Config.Format != subtitle.SRT {
		logger.Error("unsupported export format", "format", groupConfig.Config.Format.String())
	}

	if err := subtitleEntry.Update().SetIsProcessing(true).Exec(ctx); err != nil {
		logger.Error("cannot update subtitle processing status", "err", err)
		return
	}

	defer func() {
		if err := subtitleEntry.Update().SetIsProcessing(false).Exec(ctx); err != nil {
			logger.Error("cannot update subtitle processing status", "err", err)
		}
	}()

	sub := subtitle.Subtitle{
		Metadata: subtitle.Metadata{
			Language: subtitleEntry.Language,
		},
		Cues: make([]subtitle.Cue, 0),
	}

	cueEntries, err := subtitleEntry.QueryCues().All(ctx)
	if err != nil {
		logger.Error("cannot get subtitle cues", "err", err)
		return
	}

	for _, cueEntry := range cueEntries {
		cue := subtitle.Cue{
			Timestamp: subtitle.CueTimestamp{
				Start: cueEntry.TimestampStart,
				End:   cueEntry.TimestampEnd,
			},
			Content: make([]subtitle.CueContentSegment, 0),
		}

		subtitleCueContentSegmentEntries, err := cueEntry.QueryContentSegments().All(ctx)
		if err != nil {
			logger.Error("cannot get subtitle cue content segments from database", "err", err)
			return
		}

		for _, segmentEntry := range subtitleCueContentSegmentEntries {
			cue.Content = append(cue.Content, subtitle.CueContentSegment{
				Text: segmentEntry.Text,
			})
		}

		sub.Cues = append(sub.Cues, cue)
	}

	videoEntry, err := subtitleEntry.QueryVideo().Only(ctx)
	if err != nil {
		logger.Error("cannot get subtitle video", "err", err)
	}

	title := strings.ReplaceAll(strings.ReplaceAll(subtitleEntry.Title, "/", "|"), "\\", "|")

	exportFilepath := filepath.Join(
		filepath.Dir(videoEntry.Filepath),
		strings.TrimSuffix(
			filepath.Base(videoEntry.Filepath),
			filepath.Ext(videoEntry.Filepath),
		)+"."+title+"."+groupConfig.Config.Format.FileExt(),
	)

	if _, err := os.Stat(exportFilepath); err == nil {
		logger.Error("subtitle file already exist with same filepath; skipping")
		return
	} else if !os.IsNotExist(err) {
		logger.Error("cannot check if subtitle file already exist", "err", err)
		return
	}

	file, err := os.OpenFile(exportFilepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("cannot open export file for writing", "err", err)
		return
	}

	defer file.Close()

	if err := sub.Write(srt.NewWriter(file)); err != nil {
		logger.Error("cannot write subtitle to file", "err", err)
		return
	}

	if err := subtitleEntry.Update().
		SetStage(subtitleschema.StageExported).
		SetExportFormat(subtitle.SRT).
		SetExportPath(exportFilepath).
		Exec(ctx); err != nil {
		logger.Error("cannot update subtitle stage to exported", "err", err)
		return
	}
}

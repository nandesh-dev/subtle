package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/jobschema"
	"github.com/nandesh-dev/subtle/internal/jobs/export"
	"github.com/nandesh-dev/subtle/internal/jobs/extract"
	"github.com/nandesh-dev/subtle/internal/jobs/format"
	"github.com/nandesh-dev/subtle/internal/jobs/scan"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
)

var databaseJobCodes = []jobschema.Code{jobschema.CodeScan, jobschema.CodeExtract, jobschema.CodeFormat, jobschema.CodeExport}

func SetupDatabase(db *ent.Client) error {
	for _, jobCode := range databaseJobCodes {
		count, err := db.JobSchema.Update().
			Where(jobschema.CodeEQ(jobCode)).
			SetIsRunning(false).
			Save(context.Background())
		if err != nil {
			return fmt.Errorf("cannot update job status in database: %w", err)
		}

		if count == 0 {
			if err := db.JobSchema.Create().
				SetCode(jobCode).
				SetLastRun(time.Now()).
				Exec(context.Background()); err != nil {
				return fmt.Errorf("cannot add job status to database: %w", err)
			}
		}
	}

	return nil
}

const (
	configUpdateInterval  = 5 * time.Second
	defaultJobRunInterval = 30 * time.Minute
)

func StartJobRunTicker(ctx context.Context, logger *slog.Logger, configFile *configuration.File, db *ent.Client) {
	config, err := configFile.Read()
	if err != nil {
		logger.Error("failed to read configuration file", "err", err)
		return
	}

	jobRunInterval := defaultJobRunInterval
	if config.Job.Setting.Interval > 0 {
		jobRunInterval = config.Job.Setting.Interval
	} else {
		logger.Warn("job interval should be more than 0 second; using default", "default_interval", defaultJobRunInterval.String())
	}

	intervalUpdateTicker := time.NewTicker(configUpdateInterval)
	defer intervalUpdateTicker.Stop()

	jobRunTicker := time.NewTicker(jobRunInterval)
	defer jobRunTicker.Stop()

	for _, jobCode := range databaseJobCodes {
		Run(jobCode, ctx, logger, configFile, db)
	}

	for {
		select {
		case <-ctx.Done():
			return

		case <-intervalUpdateTicker.C:
			newConfig, err := configFile.Read()
			if err != nil {
				logger.Warn("failed to read configuration file; ignoring updates if any", "err", err)
				continue
			}

			newJobRunInterval := newConfig.Job.Setting.Interval

			if newJobRunInterval <= 0 {
				logger.Warn("job interval should be more than 0 second; ignoring update", "new_interval", newJobRunInterval.String())
				continue
			}

			if newJobRunInterval != jobRunInterval {
				logger.Debug("updating job interval")
				jobRunInterval = newJobRunInterval
				jobRunTicker.Reset(jobRunInterval)
			}

		case <-jobRunTicker.C:
			for _, jobCode := range databaseJobCodes {
				Run(jobCode, ctx, logger, configFile, db)
			}
		}
	}
}

func Run(jobCode jobschema.Code, ctx context.Context, logger *slog.Logger, configFile *configuration.File, db *ent.Client) {
	logger = logger.With("job", jobCode)

	if err := db.JobSchema.Update().
		Where(jobschema.CodeEQ(jobCode)).
		SetIsRunning(true).
		Exec(context.Background()); err != nil {
		logger.Error("cannot update job to running", "err", err)
		return
	}

	logger.Info("running job")

	switch jobCode {
	case jobschema.CodeScan:
		scan.Run(ctx, logger, configFile, db)
	case jobschema.CodeExtract:
		extract.Run(ctx, logger, configFile, db)
	case jobschema.CodeFormat:
		format.Run(ctx, logger, configFile, db)
	case jobschema.CodeExport:
		export.Run(ctx, logger, configFile, db)
	default:
		logger.Warn("job not found")
	}

	logger.Info("job completed")

	if err := db.JobSchema.Update().
		Where(jobschema.CodeEQ(jobCode)).
		SetIsRunning(false).
		Exec(context.Background()); err != nil {
		logger.Error("cannot update job to stopped", "err", err)
		return
	}
}

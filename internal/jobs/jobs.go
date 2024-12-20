package jobs

import (
	"context"
	"time"

	"github.com/nandesh-dev/subtle/generated/ent"
	job_schema "github.com/nandesh-dev/subtle/generated/ent/job"
	"github.com/nandesh-dev/subtle/internal/jobs/export"
	"github.com/nandesh-dev/subtle/internal/jobs/extract"
	"github.com/nandesh-dev/subtle/internal/jobs/format"
	"github.com/nandesh-dev/subtle/internal/jobs/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
)

func Init(conf *config.Config, db *ent.Client) error {
	c, err := conf.Read()
	if err != nil {
		return err
	}

	logger := logging.NewManagerLogger("job")

	jobs := []struct {
		name        string
		description string
		run         func(*config.Config, *ent.Client)
	}{
		{
			name:        "Media",
			description: "Scans the media directory for new video files and extract raw subtitle streams from it.",
			run:         media.Run,
		},
		{
			name:        "Extract",
			description: "Converts the raw subtitle streams into usable text / images based subtitles.",
			run:         extract.Run,
		},
		{
			name:        "Format",
			description: "Converts the original text / image into final text applying all the formating specified.",
			run:         format.Run,
		},
		{
			name:        "Export",
			description: "Exports subtitles to files.",
			run:         export.Run,
		},
	}

	for _, job := range jobs {
		if count, err := db.Job.Update().Where(job_schema.Name(job.name)).SetDescription(job.description).SetRunning(false).Save(context.Background()); err != nil {
			logger.Error("cannot update job info to database", "err", err)
			return err
		} else if count == 0 {
			if err := db.Job.Create().SetName(job.name).SetDescription(job.description).Exec(context.Background()); err != nil {
				logger.Error("cannot add job into to database", "err", err)
				return err
			}
		}
	}

	ticker := time.NewTicker(c.Job.Delay)
	defer ticker.Stop()

	RunAll(conf, db, jobs)

	for {
		select {
		case <-ticker.C:
			RunAll(conf, db, jobs)
		}
	}
}

func RunAll(conf *config.Config, db *ent.Client, jobs []struct {
	name        string
	description string
	run         func(*config.Config, *ent.Client)
}) {
	logger := logging.NewManagerLogger("job")

	runningJobCount, err := db.Job.Query().Where(job_schema.Running(true)).Count(context.Background())
	if err != nil {
		logger.Error("cannot get running job count from database", "err", err)
		return
	}

	if runningJobCount > 0 {
		logger.Info("some job(s) are already running! skipping")
		return
	}

	for _, job := range jobs {
		logger := logger.With("name", job.name)

		if err := db.Job.Update().Where(job_schema.Name(job.name)).SetRunning(true).Exec(context.Background()); err != nil {
			logger.Error("cannot update job to running", "err", err)
			continue
		}

		logger.Info("running job")

		job.run(conf, db)

		logger.Info("job completed")

		if err := db.Job.Update().Where(job_schema.Name(job.name)).SetRunning(false).Exec(context.Background()); err != nil {
			logger.Error("cannot update job to stopped", "err", err)
			continue
		}
	}
}

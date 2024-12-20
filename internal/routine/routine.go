package routine

import (
	"context"
	"time"

	"github.com/nandesh-dev/subtle/generated/ent"
	routine_schema "github.com/nandesh-dev/subtle/generated/ent/routine"
	"github.com/nandesh-dev/subtle/internal/routine/export"
	"github.com/nandesh-dev/subtle/internal/routine/extract"
	"github.com/nandesh-dev/subtle/internal/routine/format"
	"github.com/nandesh-dev/subtle/internal/routine/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
)

func Start(conf *config.Config, db *ent.Client) error {
	c, err := conf.Read()
	if err != nil {
		return err
	}

	logger := logging.NewManagerLogger("routine")

	routines := []struct {
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

	for _, routine := range routines {
		if count, err := db.Routine.Update().Where(routine_schema.Name(routine.name)).SetDescription(routine.description).SetRunning(false).Save(context.Background()); err != nil {
			logger.Error("cannot update routine info to database", "err", err)
			return err
		} else if count == 0 {
			if err := db.Routine.Create().SetName(routine.name).SetDescription(routine.description).Exec(context.Background()); err != nil {
				logger.Error("cannot add routine into to database", "err", err)
				return err
			}
		}
	}

	ticker := time.NewTicker(c.Routine.Delay)
	defer ticker.Stop()

	run(conf, db, routines)

	for {
		select {
		case <-ticker.C:
			run(conf, db, routines)
		}
	}
}

func run(conf *config.Config, db *ent.Client, routines []struct {
	name        string
	description string
	run         func(*config.Config, *ent.Client)
}) {
	logger := logging.NewManagerLogger("routine")

	runningRoutineCount, err := db.Routine.Query().Where(routine_schema.Running(true)).Count(context.Background())
	if err != nil {
		logger.Error("cannot get running routine count from database", "err", err)
		return
	}

	if runningRoutineCount > 0 {
		logger.Info("some routine(s) are already running! skipping")
		return
	}

	for _, routine := range routines {
		logger := logger.With("name", routine.name)

		if err := db.Routine.Update().Where(routine_schema.Name(routine.name)).SetRunning(true).Exec(context.Background()); err != nil {
			logger.Error("cannot update routine to running", "err", err)
			continue
		}

		logger.Info("running routine")

		routine.run(conf, db)

		logger.Info("routine completed")

		if err := db.Routine.Update().Where(routine_schema.Name(routine.name)).SetRunning(false).Exec(context.Background()); err != nil {
			logger.Error("cannot update routine to stopped", "err", err)
			continue
		}
	}
}

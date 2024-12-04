package routine

import (
	"fmt"
	"time"

	"github.com/nandesh-dev/subtle/internal/actions"
	"github.com/nandesh-dev/subtle/internal/routine/export"
	"github.com/nandesh-dev/subtle/internal/routine/extract"
	"github.com/nandesh-dev/subtle/internal/routine/format"
	"github.com/nandesh-dev/subtle/internal/routine/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logging"
)

func Start() error {
	if err := actions.CleanupDatebase(); err != nil {
		return fmt.Errorf("Error running cleaup database action: %v", err)
	}

	ticker := time.NewTicker(config.Config().Routine.Delay)
	defer ticker.Stop()

	run()

	for {
		select {
		case <-ticker.C:
			run()
		}
	}
}

func run() {
	logger := logging.NewManagerLogger("routine")

	routines := []struct {
		name        string
		description string
		run         func()
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
		var routineEntry database.Routine
		if err := database.Database().
			Where(
				database.Routine{
					Name: routine.name,
				},
			).
			FirstOrCreate(
				&routineEntry,
				database.Routine{
					Name:        routine.name,
					Description: routine.description,
					IsRunning:   false,
				},
			).Error; err != nil {
			logger.Error("cannot get routine from database", "err", err)
			continue
		}

		if routineEntry.IsRunning {
			logger.Error("routine already running")
			continue
		}

		routineEntry.IsRunning = true
		if err := database.Database().Save(routineEntry).Error; err != nil {
			logger.Error("cannot update routine status in database", "err", err)
		}

		logger.Info("running routine", "name", routine.name)
		routine.run()
		logger.Info("routine completed", "name", routine.name)

		routineEntry.IsRunning = false
		if err := database.Database().Save(routineEntry).Error; err != nil {
			logger.Error("cannot update routine status in database", "err", err)
		}
	}
}

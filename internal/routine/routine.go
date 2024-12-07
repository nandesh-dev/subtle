package routine

import (
	"time"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/internal/routine/export"
	"github.com/nandesh-dev/subtle/internal/routine/extract"
	"github.com/nandesh-dev/subtle/internal/routine/format"
	"github.com/nandesh-dev/subtle/internal/routine/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
)

func Start(db *ent.Client) error {
	ticker := time.NewTicker(config.Config().Routine.Delay)
	defer ticker.Stop()

	run(db)

	for {
		select {
		case <-ticker.C:
			run(db)
		}
	}
}

func run(db *ent.Client) {
	logger := logging.NewManagerLogger("routine")

	routines := []struct {
		name        string
		description string
		run         func(*ent.Client)
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
		logger.Info("running routine", "name", routine.name)
		routine.run(db)
		logger.Info("routine completed", "name", routine.name)
	}
}

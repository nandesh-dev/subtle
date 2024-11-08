package routine

import (
	"fmt"
	"time"

	"github.com/nandesh-dev/subtle/internal/routine/extract"
	"github.com/nandesh-dev/subtle/internal/routine/format"
	"github.com/nandesh-dev/subtle/internal/routine/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logger"
)

func Start() {
	ticker := time.NewTicker(config.Config().Routine.Delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			run()
		}
	}
}

func run() {
	logger.Logger().Log("Routine", fmt.Sprintf("Running routines"))

	media.Run()
	extract.Run()
	format.Run()
}

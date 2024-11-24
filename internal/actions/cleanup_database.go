package actions

import (
	"fmt"

	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"gorm.io/gorm"
)

func CleanupDatebase() error {
	var routineEntries []database.Routine
	if err := database.Database().Find(&routineEntries).Error; err != nil {
		return fmt.Errorf("Error getting routine entries from database: %v", err)
	}

	for _, routineEntry := range routineEntries {
		if routineEntry.IsRunning {
			routineEntry.IsRunning = false
			if err := database.Database().Save(routineEntry).Error; err != nil {
				return fmt.Errorf("Error updating routine status in database: %v", err)
			}
		}
	}

	var subtitleEntries []database.Subtitle
	database.Database().
		Where(database.Subtitle{IsProcessing: true}).
		Preload("Segments").
		FindInBatches(&subtitleEntries, 10, func(tx *gorm.DB, batch int) error {
			for _, subtitleEntry := range subtitleEntries {
				logger.Logger().Log("Cleanup Database Action", fmt.Sprintf("Found broken subtitle: %v; Fixing it", subtitleEntry.ID))
				if subtitleEntry.IsExtracted {
					for i := range subtitleEntry.Segments {
						subtitleEntry.Segments[i].Text = ""
					}
				} else {
					subtitleEntry.Segments = make([]database.Segment, 0)
				}

				subtitleEntry.IsProcessing = false

				database.Database().Save(subtitleEntry)
			}

			return nil
		})

	return nil
}

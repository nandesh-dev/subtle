package media

import (
	"fmt"

	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func Run() error {
	for _, rootDirectoryConfig := range config.Config().Media.RootDirectories {
		dir, _, err := filemanager.ReadDirectory(rootDirectoryConfig.Path)
		if err != nil {
			return fmt.Errorf("Error reading root directory: %v", err)
		}

		if err := syncDirectoryVideos(dir); err != nil {
			return fmt.Errorf("Error syncing directory: %v", err)
		}
	}

	return nil
}

func syncDirectoryVideos(dir *filemanager.Directory) error {
	for _, video := range dir.VideoFiles() {
		if err := db.DB().Where(db.Video{DirectoryPath: video.DirectoryPath(), Filename: video.Filename()}).FirstOrCreate(&db.Video{}, db.Video{
			DirectoryPath: video.DirectoryPath(),
			Filename:      video.Filename(),
		}).Error; err != nil {
			return fmt.Errorf("Error creating video entry: %v", err)
		}
	}

	for _, dir := range dir.Children() {
		if err := syncDirectoryVideos(&dir); err != nil {
			return err
		}
	}

	return nil
}

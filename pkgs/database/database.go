package database

import (
	"fmt"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

type Video struct {
	ID int `gorm:"primaryKey"`

	DirectoryPath string
	Filename      string
	Subtitles     []Subtitle `gorm:"foreignKey:VideoID"`
}

type Subtitle struct {
	ID      int `gorm:"primaryKey"`
	VideoID int

	Title        string
	Language     string
	IsProcessing bool
	Segments     []Segment `gorm:"foreignKey:SubtitleID"`

	IsExtracted bool
	IsFormated  bool
	IsExported  bool

	ImportIsExternal       bool
	ImportFormat           string
	ImportVideoStreamIndex int

	ExportPath   string
	ExportFormat string
}

type Segment struct {
	ID         int `gorm:"primaryKey"`
	SubtitleID int

	StartTime time.Duration
	EndTime   time.Duration
	Text      string

	OriginalText  string
	OriginalImage []byte
}

type Routine struct {
	Name        string `gorm:"primaryKey"`
	Description string
	IsRunning   bool
}

func Database() *gorm.DB {
	return database
}

func Init() error {
	db, err := gorm.Open(sqlite.Open(config.Config().Database.Path), &gorm.Config{
		FullSaveAssociations: true,
	})
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}

	if err = db.AutoMigrate(&Video{}, &Subtitle{}, &Segment{}, &Routine{}); err != nil {
		return fmt.Errorf("Error auto migrating database: %v", err)
	}

	database = db
	return nil
}

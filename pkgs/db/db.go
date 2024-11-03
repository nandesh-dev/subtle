package db

import (
	"fmt"
	"time"

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

	Title    string
	Language string
	Segments []Segment `gorm:"foreignKey:SubtitleID"`

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

func DB() *gorm.DB {
	return database
}

func Init() error {
	db, err := gorm.Open(sqlite.Open("config/database.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}

	if err = db.AutoMigrate(&Video{}, &Subtitle{}, &Segment{}); err != nil {
		return fmt.Errorf("Error auto migrating database: %v", err)
	}

	database = db
	return nil
}

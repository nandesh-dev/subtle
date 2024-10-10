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

	Language string
	Filepath string
	IsImage  bool
	Segments []Segment `gorm:"foreignKey:SubtitleID"`
}

type Segment struct {
	ID         int `gorm:"primaryKey"`
	SubtitleID int

	StartTime time.Duration
	EndTime   time.Duration
	ImageData []byte
	Text      string
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

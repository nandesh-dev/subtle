package configuration

import (
	"fmt"
	"os"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/language"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"gopkg.in/yaml.v3"
)

type ScanningGroup struct {
	DirectoryPath string `yaml:"directory_path"`
}

type ExtractingGroup struct {
	Condition ExtractingCondition `yaml:"condition"`
	Priority  ExtractingPriority  `yaml:"priority"`
	Limit     int                 `yaml:"limit"`
}

type ExtractingCondition struct {
	Formats   []subtitle.Format `yaml:"formats"`
	Languages []language.Tag    `yaml:"languages"`
}

type ExtractingPriority struct {
	Format       map[subtitle.Format]int `yaml:"format"`
	Language     map[language.Tag]int    `yaml:"language"`
	TitleKeyword map[string]int          `yaml:"title_keyword"`
}

type FormatingGroup struct {
	Condition FormatingCondition `yaml:"condition"`
	Priority  FormatingPriority  `yaml:"priority"`
	Config    FormatingConfig    `yaml:"config"`
	Limit     int                `yaml:"limit"`
}

type FormatingCondition struct {
	Formats   []subtitle.Format `yaml:"formats"`
	Languages []language.Tag    `yaml:"languages"`
}

type FormatingPriority struct {
	Format       map[subtitle.Format]int `yaml:"format"`
	Language     map[language.Tag]int    `yaml:"language"`
	TitleKeyword map[string]int          `yaml:"title_keyword"`
}

type FormatingConfig struct {
	WordMappings []FormatingConfigWordMapping `yaml:"word_mappings"`
}

type FormatingConfigWordMapping struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type ExportingGroup struct {
	Condition ExportingCondition `yaml:"condition"`
	Priority  ExtractingPriority `yaml:"priority"`
	Config    ExportingConfig    `yaml:"config"`
	Limit     int                `yaml:"limit"`
}

type ExportingCondition struct {
	Formats   []subtitle.Format `yaml:"formats"`
	Languages []language.Tag    `yaml:"languages"`
}

type ExportingPriority struct {
	Format       map[subtitle.Format]int `yaml:"format"`
	Language     map[language.Tag]int    `yaml:"language"`
	TitleKeyword map[string]int          `yaml:"title_keyword"`
}

type ExportingConfig struct {
	Format subtitle.Format `yaml:"format"`
}

type Job struct {
	Setting    JobSetting        `yaml:"job_setting"`
	Scanning   []ScanningGroup   `yaml:"scanning"`
	Extracting []ExtractingGroup `yaml:"extracting"`
	Formating  []FormatingGroup  `yaml:"formating"`
	Exporting  []ExportingGroup  `yaml:"exporting"`
}

type JobSetting struct {
	Interval time.Duration `yaml:"interval"`
}

type Config struct {
	Job Job `yaml:"job"`
}

type File struct {
	path string
}

var Default Config = Config{
	Job: Job{
		Scanning:   []ScanningGroup{},
		Extracting: []ExtractingGroup{},
		Formating:  []FormatingGroup{},
		Exporting:  []ExportingGroup{},
		Setting: JobSetting{
			Interval: 15 * time.Minute,
		},
	},
}

func Open(path string) (*File, error) {
	file := File{path: path}

	if _, err := file.Read(); os.IsNotExist(err) {
		file.Write(Default)
	} else if err != nil {
		return nil, err
	}

	return &file, nil
}

func (f File) Read() (*Config, error) {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (f File) Write(config Config) error {
	output, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("Error marshaling file: %w", err)
	}

	if err := os.WriteFile(f.path, output, 0644); err != nil {
		return fmt.Errorf("Error writing config: %w", err)
	}

	return nil
}

func (f File) ReadString() (string, error) {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

  return string(file), nil
}

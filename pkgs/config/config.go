package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type MediaDirectory struct {
	Path       string
	Extraction Extraction
	Formating  Formating
	Exporting  Exporting
}

type Extraction struct {
	Enable                 bool
	RawStreamTitleKeywords []string `yaml:"raw_stream_title_keywords"`
	Formats                Formats
}

type Formats struct {
	ASS ASS
	PGS PGS
}

type ASS struct {
	Enable    bool
	Languages []language.Tag
}

type PGS struct {
	Enable    bool
	Languages []language.Tag
}

type Formating struct {
	Enable             bool
	TextBasedSubtitle  TextBasedSubtitle  `yaml:"text_based_subtitle"`
	ImageBasedSubtitle ImageBasedSubtitle `yaml:"image_based_subtitle"`
}

type TextBasedSubtitle struct {
	CharactorMappings []CharactorMapping `yaml:"charactor_mappings"`
}

type ImageBasedSubtitle struct {
	CharactorMappings []CharactorMapping `yaml:"charactor_mappings"`
}

type CharactorMapping struct {
	Language language.Tag
	Mappings []Mapping
}

type Mapping struct {
	From string
	To   string
}

type Exporting struct {
	Enable bool
	Format string
}

type Server struct {
	Web      Web
	Database Database
	Routine  Routine
	Logging  Logging
}

type Web struct {
	Port                 int
	COROrigins           []string `yaml:"cor_origins"`
	EnableGRPCReflection bool     `yaml:"enable_grpc_reflection"`
	ServeDirectory       string   `yaml:"serve_directory"`
}

type Database struct {
	Path string
}

type Routine struct {
	Delay time.Duration
}

type Logging struct {
	Path string
}

type Data struct {
	MediaDirectories []MediaDirectory `yaml:"watch_directories"`
	Server
}

type Config struct {
	Path string
}

func Open(path string) (*Config, error) {
	basepath := filepath.Dir(path)

	defaultConfig := Data{
		Server: Server{
			Web: Web{
				Port:                 3000,
				COROrigins:           make([]string, 0),
				EnableGRPCReflection: false,
				ServeDirectory:       "/public",
			},
			Database: Database{
				Path: filepath.Join(basepath, "database.db"),
			},
			Routine: Routine{
				Delay: time.Minute * 15,
			},
			Logging: Logging{
				Path: filepath.Join(basepath, "logs.log"),
			},
		},
		MediaDirectories: []MediaDirectory{
			{
				Path: "/media",
				Extraction: Extraction{
					Enable:                 false,
					RawStreamTitleKeywords: []string{"Full", "Dialogue"},
					Formats: Formats{
						ASS: ASS{
							Enable:    true,
							Languages: []language.Tag{language.English},
						},
						PGS: PGS{
							Enable:    true,
							Languages: []language.Tag{language.English},
						},
					},
				},
				Formating: Formating{
					Enable: false,
					TextBasedSubtitle: TextBasedSubtitle{
						CharactorMappings: []CharactorMapping{},
					},
					ImageBasedSubtitle: ImageBasedSubtitle{
						CharactorMappings: []CharactorMapping{
							{
								Language: language.English,
								Mappings: []Mapping{
									{
										From: "|",
										To:   "I",
									},
								},
							},
						},
					},
				},
				Exporting: Exporting{
					Enable: false,
					Format: "srt",
				},
			},
		},
	}

	config := Config{
		Path: path,
	}

	if _, err := config.Read(); os.IsNotExist(err) {
		config.Write(defaultConfig)
	} else if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Read() (*Data, error) {
	file, err := os.ReadFile(c.Path)
	if err != nil {
		return nil, err
	}

	var data Data
	if err := yaml.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Config) Write(data Data) error {
	output, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("Error marshaling file: %w", err)
	}

	if err := os.WriteFile(c.Path, output, 0644); err != nil {
		return fmt.Errorf("Error writing config: %w", err)
	}

	return nil
}

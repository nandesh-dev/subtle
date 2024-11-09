package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
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

type t struct {
	MediaDirectories []MediaDirectory `yaml:"watch_directories"`
	Server
}

var (
	config t
	path   string
	once   sync.Once
)

func Config() *t {
	return &config
}

func Init(basepath string) (e error) {
	once.Do(func() {
		config = t{
			Server: Server{
				Web: Web{
					Port:                 3000,
					COROrigins:           make([]string, 0),
					EnableGRPCReflection: false,
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
						Format: "srt",
					},
				},
			},
		}

		path = filepath.Join(basepath, "config.yaml")

		file, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				file, err := os.Create(filepath.Join(basepath, "config.yaml"))
				if err != nil {
					e = fmt.Errorf("Error creating config file: %v", err)
					return
				}
				file.Close()

				e = Write()
				return
			}

			e = fmt.Errorf("Error reading config file: %v", err)
			return
		}

		if err := yaml.Unmarshal(file, &config); err != nil {
			e = fmt.Errorf("Error unmarshaling file: %v", err)
		}
	})

	return
}

func Write() error {
	if path == "" {
		return fmt.Errorf("Config not initilized")
	}

	output, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("Error marshaling file: %v", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("Error writing config: %v", err)
	}

	return nil
}

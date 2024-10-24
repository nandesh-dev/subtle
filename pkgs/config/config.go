package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Port              int
	COROrigins        []string
	GRPCReflection    bool
	DatabaseDirectory string
}

type AutoExtractFormat struct {
	ASS AutoExtractASS
}

type AutoExtractASS struct {
	Enabled bool
}

type AutoExtract struct {
	Languages              []language.Tag
	Formats                AutoExtractFormat
	RawStreamTitleKeywords []string
}

type RootDirectory struct {
	Path        string
	AutoExtract AutoExtract
}

type Media struct {
	RootDirectories []RootDirectory
}

type t struct {
	Server Server
	Media  Media
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
				Port:              3000,
				GRPCReflection:    false,
				DatabaseDirectory: filepath.Join(basepath, "db"),
			},
			Media: Media{
				RootDirectories: []RootDirectory{
					{
						Path: "/media",
						AutoExtract: AutoExtract{
							Languages: []language.Tag{language.English},
							Formats: AutoExtractFormat{
								ASS: AutoExtractASS{
									Enabled: true,
								},
							},
							RawStreamTitleKeywords: []string{"Full", "Dialogue"},
						},
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

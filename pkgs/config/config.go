package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Port           int
	GRPCReflection bool
}

type AutoExtract struct {
	Formats      []string
	Languages    []string
	OutputFormat string
}

type WatchDirectory struct {
	Path        string
	AutoExtract AutoExtract
}

type Media struct {
	WatchDirectories []WatchDirectory
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
				Port:           3000,
				GRPCReflection: false,
			},
			Media: Media{
				WatchDirectories: []WatchDirectory{
					{
						Path: "/media",
						AutoExtract: AutoExtract{
							Languages:    []string{"eng"},
							Formats:      []string{"ass"},
							OutputFormat: "srt",
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

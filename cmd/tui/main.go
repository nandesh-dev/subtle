package main

import (
	"fmt"
	"log"

	"github.com/nandesh-dev/subtle/internal/filemanager"
)

func main() {
	directory, err := filemanager.ReadDirectory(".")
	if err != nil {
		log.Fatal(err)
	}

	stats, err := directory.Videos[0].Stats()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Stats: %v\n", stats)
}

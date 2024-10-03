package main

import (
	"log"

	"github.com/nandesh-dev/subtle/internal/routine/library"
	"github.com/nandesh-dev/subtle/pkgs/config"
)

func main() {
	if err := config.Init("/config"); err != nil {
		log.Fatal(err)
	}

	library.RunLibraryRoutine()
}

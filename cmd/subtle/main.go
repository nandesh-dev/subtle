package main

import (
	"fmt"
	"log"

	"github.com/nandesh-dev/subtle/internal/routine/extract"
	"github.com/nandesh-dev/subtle/internal/routine/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/db"
)

func main() {
	fmt.Println("Initilizing config")
	if err := config.Init("./config"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initilizing database")
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running media routine")
	media.Run()

	fmt.Println("Running extract routine")
	warns := extract.Run()

	for _, warning := range warns.Warnings() {
		fmt.Println(warning)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/nandesh-dev/subtle/internal/server"
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

	server := server.New()
	server.Listen(3000, true)
}

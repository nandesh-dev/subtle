package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nandesh-dev/subtle/internal/filemanager"
	"github.com/nandesh-dev/subtle/internal/pb/library"
	"github.com/nandesh-dev/subtle/internal/pgs"
	"github.com/nandesh-dev/subtle/internal/services"
	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/tesseract"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dir, err := filemanager.ReadDirectory("/media")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dir)

	videos, warnings := dir.VideoFiles()
	fmt.Println(warnings)

	rawStreams, err := subtitle.ExtractRawStreams(&videos[0])
	if err != nil {
		log.Fatal(err)
	}

	subtitle, _, err := pgs.DecodeSubtitle(&rawStreams[0])
	if err != nil {
		log.Fatal(err)
	}

	tes := tesseract.NewClient()

	for _, segment := range subtitle.Segments() {
		for _, img := range segment.Images() {
			text, err := tes.ExtractTextFromImage(img, language.English)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(text)
		}
	}

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln("failed to create listener: ", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	libraryService := services.Services().Library
	library.RegisterLibraryServiceServer(s, &libraryService)

	if err := s.Serve(listener); err != nil {
		log.Fatalln("failed to serve: ", err)
	}
}

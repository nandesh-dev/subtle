package services

import "github.com/nandesh-dev/subtle/internal/services/library"

type services struct {
	Library library.LibraryServiceServer
}

func Services() services {
	return services{
		Library: library.LibraryServiceServer{},
	}
}

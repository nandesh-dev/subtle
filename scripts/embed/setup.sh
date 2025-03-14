mkdir -p generated/embed

cat > "generated/embed/embed.go" <<EOF
// GENERATED CODE! DO NOT EDIT!

package embed

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed frontend/*
var embeddedFrontendFiles embed.FS

func GetFrontendFilesystem() fs.FS {
	frontendFS, err := fs.Sub(embeddedFrontendFiles, "frontend")
	if err != nil {
		log.Panicln("Failed to create a sub filesystem! Was the frontend files missing during build time?")
	}

	return frontendFS
}
EOF


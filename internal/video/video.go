package video

type VideoFile struct {
	Path   string
	Format VideoFileFileFormat
}

type VideoFileFileFormat int

const (
	MP4 VideoFileFileFormat = iota
	MKV
	AVI
	MOV
)

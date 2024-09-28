package video

type File struct {
	Path   string
	Format Format
}

type Format int

const (
	MP4 Format = iota
	MKV
	AVI
	MOV
)

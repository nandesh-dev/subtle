package subtitle

type File struct {
	Path   string
	Format Format
}

type Format int

const (
	ASS Format = iota
	PGS
)

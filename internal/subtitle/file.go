package subtitle

type File struct {
	Path   string
	Format Format
}

type Format int

const (
	SRT Format = iota
	ASS
	SSA
	IDX
	SUB
	PGS
)

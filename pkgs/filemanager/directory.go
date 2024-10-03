package filemanager

type Directory struct {
	path           string
	children       []Directory
	videos         []VideoFile
	extraSubtitles []SubtitleFile
}

func NewDirectory(path string) *Directory {
	return &Directory{
		path:           path,
		children:       make([]Directory, 0),
		videos:         make([]VideoFile, 0),
		extraSubtitles: make([]SubtitleFile, 0),
	}
}

func (d *Directory) AddVideoFile(file VideoFile) {
	d.videos = append(d.videos, file)
}

func (d *Directory) AddExtraSubtitleFile(file SubtitleFile) {
	d.extraSubtitles = append(d.extraSubtitles, file)
}

func (d *Directory) AddChild(child Directory) {
	d.children = append(d.children, child)
}

func (d *Directory) VideoFiles() []VideoFile {
	return d.videos
}

func (d *Directory) ExtraSubtitleFiles() []SubtitleFile {
	return d.extraSubtitles
}

func (d *Directory) Children() []Directory {
	return d.children
}

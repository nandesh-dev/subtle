package filemanager

type Directory struct {
	path      string
	children  []Directory
	videos    []VideoFile
	subtitles []SubtitleFile
}

func (d *Directory) VideoFiles() []VideoFile {
	return d.videos
}

func (d *Directory) SubtitleFiles() []SubtitleFile {
	return d.subtitles
}

func (d *Directory) Children() []Directory {
	return d.children
}

func (d *Directory) Path() string {
	return d.path
}

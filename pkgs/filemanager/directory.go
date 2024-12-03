package filemanager

type Directory struct {
	Path          string
	ChildrenPaths []string
	Videos        []VideoFile
	Subtitles     []SubtitleFile
}

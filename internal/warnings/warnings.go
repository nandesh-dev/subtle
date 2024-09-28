package warnings

type WarningList struct {
	Warnings []error
}

func NewWarningList() *WarningList {
	return &WarningList{
		Warnings: make([]error, 0),
	}
}

func (w *WarningList) AddWarning(e error) {
	w.Warnings = append(w.Warnings, e)
}

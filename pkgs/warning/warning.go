package warning

type WarningList struct {
	warnings []error
}

func NewWarningList() *WarningList {
	return &WarningList{
		warnings: make([]error, 0),
	}
}

func (w *WarningList) AddWarning(e error) {
	w.warnings = append(w.warnings, e)
}

func (w *WarningList) Warnings() []error {
	return w.warnings
}

func (w *WarningList) Append(warnings WarningList) {
	w.warnings = append(w.warnings, warnings.warnings...)
}

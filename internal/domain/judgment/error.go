package judgment

import "errors"

var (
	ErrJudgmentNotFound            = errors.New("judgment not found")
	ErrEmptyPrescription           = errors.New("empty prescription")
	ErrInvalidPrescribedMenuLabels = errors.New("invalid prescribed menu labels")
)

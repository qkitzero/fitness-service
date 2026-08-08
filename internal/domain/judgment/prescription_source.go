package judgment

type PrescriptionSource string

const (
	PrescriptionSourceElement   PrescriptionSource = "element"
	PrescriptionSourceFixed     PrescriptionSource = "fixed"
	PrescriptionSourceAgeDecade PrescriptionSource = "age_decade"
	PrescriptionSourceManual    PrescriptionSource = "manual"
)

func (p PrescriptionSource) String() string {
	return string(p)
}

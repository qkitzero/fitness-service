package measurementitem

import "fmt"

type Unit string

const (
	UnitKg      Unit = "kg"
	UnitCm      Unit = "cm"
	UnitSec     Unit = "sec"
	UnitCount   Unit = "count"
	UnitMmHg    Unit = "mmHg"
	UnitPercent Unit = "percent"
	UnitBpm     Unit = "bpm"
	UnitLevel   Unit = "level"
)

func (u Unit) String() string {
	return string(u)
}

func NewUnit(s string) (Unit, error) {
	switch Unit(s) {
	case UnitKg, UnitCm, UnitSec, UnitCount, UnitMmHg, UnitPercent, UnitBpm, UnitLevel:
		return Unit(s), nil
	default:
		return Unit(""), fmt.Errorf("invalid unit")
	}
}

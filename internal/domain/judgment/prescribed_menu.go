package judgment

import (
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type PrescribedMenu interface {
	Source() PrescriptionSource
	Element() *measurementitem.Element
	Part() *training.Part
	TrainingMenuID() training.TrainingMenuID
	TrainingMenuName() training.Name
	Amount() training.Amount
	Unit() training.Unit
	Sets() training.Sets
}

type prescribedMenu struct {
	source           PrescriptionSource
	element          *measurementitem.Element
	part             *training.Part
	trainingMenuID   training.TrainingMenuID
	trainingMenuName training.Name
	amount           training.Amount
	unit             training.Unit
	sets             training.Sets
}

func (p prescribedMenu) Source() PrescriptionSource {
	return p.source
}

func (p prescribedMenu) Element() *measurementitem.Element {
	if p.element == nil {
		return nil
	}
	e := *p.element
	return &e
}

func (p prescribedMenu) Part() *training.Part {
	if p.part == nil {
		return nil
	}
	part := *p.part
	return &part
}

func (p prescribedMenu) TrainingMenuID() training.TrainingMenuID {
	return p.trainingMenuID
}

func (p prescribedMenu) TrainingMenuName() training.Name {
	return p.trainingMenuName
}

func (p prescribedMenu) Amount() training.Amount {
	return p.amount
}

func (p prescribedMenu) Unit() training.Unit {
	return p.unit
}

func (p prescribedMenu) Sets() training.Sets {
	return p.sets
}

func newPrescribedMenu(
	source PrescriptionSource,
	element *measurementitem.Element,
	part *training.Part,
	trainingMenuID training.TrainingMenuID,
	trainingMenuName training.Name,
	amount training.Amount,
	unit training.Unit,
	sets training.Sets,
) PrescribedMenu {
	p := &prescribedMenu{
		source:           source,
		trainingMenuID:   trainingMenuID,
		trainingMenuName: trainingMenuName,
		amount:           amount,
		unit:             unit,
		sets:             sets,
	}
	if element != nil {
		e := *element
		p.element = &e
	}
	if part != nil {
		pt := *part
		p.part = &pt
	}
	return p
}

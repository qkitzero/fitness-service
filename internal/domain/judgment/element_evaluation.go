package judgment

import (
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
)

type ElementEvaluation interface {
	Element() measurementitem.Element
	ZScore() standard.ZScore
	Rank() standard.Rank
}

type elementEvaluation struct {
	element measurementitem.Element
	zScore  standard.ZScore
	rank    standard.Rank
}

func (e elementEvaluation) Element() measurementitem.Element {
	return e.element
}

func (e elementEvaluation) ZScore() standard.ZScore {
	return e.zScore
}

func (e elementEvaluation) Rank() standard.Rank {
	return e.rank
}

func newElementEvaluation(
	element measurementitem.Element,
	zScore standard.ZScore,
	rank standard.Rank,
) ElementEvaluation {
	return &elementEvaluation{
		element: element,
		zScore:  zScore,
		rank:    rank,
	}
}

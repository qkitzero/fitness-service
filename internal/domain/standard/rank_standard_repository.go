package standard

import (
	"context"
)

type RankStandardRepository interface {
	List(ctx context.Context) ([]RankStandard, error)
}

package training

import (
	"context"
)

type PrescriptionRuleRepository interface {
	ListElementMenus(ctx context.Context) ([]ElementMenu, error)
	ListFixedMenus(ctx context.Context) ([]FixedMenu, error)
	ListAgeDecadeMenus(ctx context.Context) ([]AgeDecadeMenu, error)
}

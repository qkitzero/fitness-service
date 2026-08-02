package judgment

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type JudgmentResult struct {
	MeasurementID measurement.MeasurementID
	IsDraft       bool
	Evaluation    judgment.Evaluation
	Advice        *judgment.Advice
}

type AdvicePatch struct {
	Advice    *judgment.Advice
	HasAdvice bool
}

type JudgmentUsecase interface {
	GetJudgment(ctx context.Context, measurementID measurement.MeasurementID) (JudgmentResult, error)
	UpsertJudgmentAdvice(ctx context.Context, measurementID measurement.MeasurementID, patch AdvicePatch) (judgment.Judgment, error)
}

type judgmentUsecase struct {
	authService          auth.AuthService
	userService          user.UserService
	judgmentRepo         judgment.JudgmentRepository
	measurementRepo      measurement.MeasurementRepository
	customerRepo         customer.CustomerRepository
	measurementItemRepo  measurementitem.MeasurementItemRepository
	ageGroupStandardRepo standard.AgeGroupStandardRepository
	rankStandardRepo     standard.RankStandardRepository
}

func NewJudgmentUsecase(
	authService auth.AuthService,
	userService user.UserService,
	judgmentRepo judgment.JudgmentRepository,
	measurementRepo measurement.MeasurementRepository,
	customerRepo customer.CustomerRepository,
	measurementItemRepo measurementitem.MeasurementItemRepository,
	ageGroupStandardRepo standard.AgeGroupStandardRepository,
	rankStandardRepo standard.RankStandardRepository,
) JudgmentUsecase {
	return &judgmentUsecase{
		authService:          authService,
		userService:          userService,
		judgmentRepo:         judgmentRepo,
		measurementRepo:      measurementRepo,
		customerRepo:         customerRepo,
		measurementItemRepo:  measurementItemRepo,
		ageGroupStandardRepo: ageGroupStandardRepo,
		rankStandardRepo:     rankStandardRepo,
	}
}

func (u *judgmentUsecase) verifyTenantMembership(ctx context.Context, tenantID tenant.TenantID) error {
	groupIDs, err := u.userService.ListMyGroups(ctx)
	if err != nil {
		return err
	}

	for _, id := range groupIDs {
		if id == tenantID.String() {
			return nil
		}
	}

	return tenant.ErrNotMember
}

func (u *judgmentUsecase) findOwnedMeasurement(ctx context.Context, measurementID measurement.MeasurementID) (measurement.Measurement, customer.Customer, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, nil, err
	}

	foundMeasurement, err := u.measurementRepo.FindByID(ctx, measurementID)
	if err != nil {
		return nil, nil, err
	}

	foundCustomer, err := u.customerRepo.FindByID(ctx, foundMeasurement.CustomerID())
	if err != nil {
		return nil, nil, err
	}

	if err := u.verifyTenantMembership(ctx, foundCustomer.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, nil, judgment.ErrJudgmentNotFound
		}
		return nil, nil, err
	}

	return foundMeasurement, foundCustomer, nil
}

func (u *judgmentUsecase) findJudgment(ctx context.Context, measurementID measurement.MeasurementID) (judgment.Judgment, error) {
	foundJudgment, err := u.judgmentRepo.FindByMeasurementID(ctx, measurementID)
	if err != nil && !errors.Is(err, judgment.ErrJudgmentNotFound) {
		return nil, err
	}

	return foundJudgment, nil
}

func (u *judgmentUsecase) evaluate(ctx context.Context, foundMeasurement measurement.Measurement, foundCustomer customer.Customer) (judgment.Evaluation, error) {
	age := foundMeasurement.AgeAtMeasurement().Int()

	gender, err := standard.NewGender(foundCustomer.Gender().String())
	if err != nil {
		log.Printf("GetJudgment: measurement %s: no standards exist for customer gender %q, returning an empty evaluation", foundMeasurement.ID(), foundCustomer.Gender())
		return judgment.NewEvaluation(foundMeasurement, nil, gender, age, nil, nil), nil
	}

	entries := foundMeasurement.Entries()
	measurementItemIDs := make([]measurementitem.MeasurementItemID, 0, len(entries))
	for _, entry := range entries {
		measurementItemIDs = append(measurementItemIDs, entry.MeasurementItemID())
	}

	measurementItems, err := u.measurementItemRepo.FindByIDs(ctx, measurementItemIDs)
	if err != nil {
		return nil, err
	}

	ageGroupStandards, err := u.ageGroupStandardRepo.ListByItemIDsAndGender(ctx, measurementItemIDs, gender)
	if err != nil {
		return nil, err
	}

	rankStandards, err := u.rankStandardRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(entries) > 0 && len(ageGroupStandards) == 0 {
		log.Printf("GetJudgment: measurement %s: no age group standards are registered for gender %q, returning an empty evaluation", foundMeasurement.ID(), gender)
	}

	return judgment.NewEvaluation(foundMeasurement, measurementItems, gender, age, ageGroupStandards, rankStandards), nil
}

func (u *judgmentUsecase) GetJudgment(ctx context.Context, measurementID measurement.MeasurementID) (JudgmentResult, error) {
	foundMeasurement, foundCustomer, err := u.findOwnedMeasurement(ctx, measurementID)
	if err != nil {
		return JudgmentResult{}, err
	}

	evaluation, err := u.evaluate(ctx, foundMeasurement, foundCustomer)
	if err != nil {
		return JudgmentResult{}, err
	}

	foundJudgment, err := u.findJudgment(ctx, measurementID)
	if err != nil {
		return JudgmentResult{}, err
	}

	result := JudgmentResult{
		MeasurementID: measurementID,
		IsDraft:       foundMeasurement.IsDraft(),
		Evaluation:    evaluation,
	}
	if foundJudgment != nil {
		result.Advice = foundJudgment.Advice()
	}

	return result, nil
}

func (u *judgmentUsecase) UpsertJudgmentAdvice(ctx context.Context, measurementID measurement.MeasurementID, patch AdvicePatch) (judgment.Judgment, error) {
	if _, _, err := u.findOwnedMeasurement(ctx, measurementID); err != nil {
		return nil, err
	}

	foundJudgment, err := u.findJudgment(ctx, measurementID)
	if err != nil {
		return nil, err
	}

	if foundJudgment != nil {
		if !patch.HasAdvice {
			return foundJudgment, nil
		}

		foundJudgment.Update(patch.Advice)

		if err := u.judgmentRepo.Upsert(ctx, foundJudgment); err != nil {
			return nil, err
		}

		return foundJudgment, nil
	}

	now := time.Now().UTC()
	createdJudgment := judgment.NewJudgment(measurementID, patch.Advice, now, now)

	if createdJudgment.Advice() == nil {
		return createdJudgment, nil
	}

	if err := u.judgmentRepo.Upsert(ctx, createdJudgment); err != nil {
		return nil, err
	}

	return createdJudgment, nil
}

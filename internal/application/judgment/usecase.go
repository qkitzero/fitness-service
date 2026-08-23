package judgment

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/judgment"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
	"github.com/qkitzero/fitness-service/internal/domain/standard"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
	"github.com/qkitzero/fitness-service/internal/domain/training"
)

type JudgmentResult struct {
	MeasurementID measurement.MeasurementID
	IsDraft       bool
	Evaluation    judgment.Evaluation
	Prescription  judgment.Prescription
	Advice        *judgment.Advice
}

type OrganizationJudgmentResult struct {
	CustomerID       customer.CustomerID
	MeasurementID    measurement.MeasurementID
	MeasuredOn       measurement.MeasuredOn
	AgeAtMeasurement measurement.AgeAtMeasurement
	IsDraft          bool
	Evaluation       judgment.Evaluation
}

type AdvicePatch struct {
	Advice    *judgment.Advice
	HasAdvice bool
}

type PrescribedMenuInput struct {
	Element        *measurementitem.Element
	Part           *training.Part
	TrainingMenuID training.TrainingMenuID
	Amount         training.Amount
	Unit           training.Unit
	Sets           training.Sets
}

type JudgmentUsecase interface {
	GetJudgment(ctx context.Context, measurementID measurement.MeasurementID) (JudgmentResult, error)
	ListOrganizationJudgments(ctx context.Context, organizationID organization.OrganizationID, includeInactive bool) ([]OrganizationJudgmentResult, error)
	UpsertJudgmentAdvice(ctx context.Context, measurementID measurement.MeasurementID, patch AdvicePatch) (judgment.Judgment, error)
	UpsertPrescription(ctx context.Context, measurementID measurement.MeasurementID, menus []PrescribedMenuInput) (judgment.Prescription, error)
	DeletePrescription(ctx context.Context, measurementID measurement.MeasurementID) error
}

type judgmentUsecase struct {
	authService                auth.AuthService
	userService                user.UserService
	judgmentRepo               judgment.JudgmentRepository
	prescribedMenuOverrideRepo judgment.PrescribedMenuOverrideRepository
	measurementRepo            measurement.MeasurementRepository
	customerRepo               customer.CustomerRepository
	organizationRepo           organization.OrganizationRepository
	measurementItemRepo        measurementitem.MeasurementItemRepository
	ageGroupStandardRepo       standard.AgeGroupStandardRepository
	rankStandardRepo           standard.RankStandardRepository
	trainingMenuRepo           training.TrainingMenuRepository
	prescriptionRuleRepo       training.PrescriptionRuleRepository
}

func NewJudgmentUsecase(
	authService auth.AuthService,
	userService user.UserService,
	judgmentRepo judgment.JudgmentRepository,
	prescribedMenuOverrideRepo judgment.PrescribedMenuOverrideRepository,
	measurementRepo measurement.MeasurementRepository,
	customerRepo customer.CustomerRepository,
	organizationRepo organization.OrganizationRepository,
	measurementItemRepo measurementitem.MeasurementItemRepository,
	ageGroupStandardRepo standard.AgeGroupStandardRepository,
	rankStandardRepo standard.RankStandardRepository,
	trainingMenuRepo training.TrainingMenuRepository,
	prescriptionRuleRepo training.PrescriptionRuleRepository,
) JudgmentUsecase {
	return &judgmentUsecase{
		authService:                authService,
		userService:                userService,
		judgmentRepo:               judgmentRepo,
		prescribedMenuOverrideRepo: prescribedMenuOverrideRepo,
		measurementRepo:            measurementRepo,
		customerRepo:               customerRepo,
		organizationRepo:           organizationRepo,
		measurementItemRepo:        measurementItemRepo,
		ageGroupStandardRepo:       ageGroupStandardRepo,
		rankStandardRepo:           rankStandardRepo,
		trainingMenuRepo:           trainingMenuRepo,
		prescriptionRuleRepo:       prescriptionRuleRepo,
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

func itemCodesWithoutAgeGroupStandards(items []measurementitem.MeasurementItem, ageGroupStandards []standard.AgeGroupStandard) []string {
	registeredItemIDs := make(map[measurementitem.MeasurementItemID]struct{}, len(ageGroupStandards))
	for _, ageGroupStandard := range ageGroupStandards {
		registeredItemIDs[ageGroupStandard.MeasurementItemID()] = struct{}{}
	}

	codes := make([]string, 0, len(items))
	for _, item := range items {
		if item.ScoreDirection() == nil {
			continue
		}
		if _, ok := registeredItemIDs[item.ID()]; ok {
			continue
		}
		codes = append(codes, item.Code().String())
	}
	sort.Strings(codes)

	return codes
}

func itemCodesOutsideAgeGroups(items []measurementitem.MeasurementItem, ageGroupStandards []standard.AgeGroupStandard, age int) []string {
	outsideByItemID := make(map[measurementitem.MeasurementItemID]bool, len(items))
	for _, ageGroupStandard := range ageGroupStandards {
		measurementItemID := ageGroupStandard.MeasurementItemID()
		if ageGroupStandard.AgeRange().Contains(age) {
			outsideByItemID[measurementItemID] = false
			continue
		}
		if _, ok := outsideByItemID[measurementItemID]; !ok {
			outsideByItemID[measurementItemID] = true
		}
	}

	codes := make([]string, 0, len(items))
	for _, item := range items {
		if !outsideByItemID[item.ID()] {
			continue
		}
		codes = append(codes, item.Code().String())
	}
	sort.Strings(codes)

	return codes
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
	} else if unregisteredCodes := itemCodesWithoutAgeGroupStandards(measurementItems, ageGroupStandards); len(unregisteredCodes) > 0 {
		log.Printf("GetJudgment: measurement %s: the age group standards of gender %q are missing for items %v", foundMeasurement.ID(), gender, unregisteredCodes)
	}

	if outsideCodes := itemCodesOutsideAgeGroups(measurementItems, ageGroupStandards, age); len(outsideCodes) > 0 {
		log.Printf("GetJudgment: measurement %s: age %d is outside the registered age groups of gender %q for items %v", foundMeasurement.ID(), age, gender, outsideCodes)
	}

	return judgment.NewEvaluation(foundMeasurement, measurementItems, gender, age, ageGroupStandards, rankStandards), nil
}

func (u *judgmentUsecase) prescribe(ctx context.Context, foundMeasurement measurement.Measurement, evaluation judgment.Evaluation) (judgment.Prescription, error) {
	overrides, err := u.prescribedMenuOverrideRepo.ListByMeasurementID(ctx, foundMeasurement.ID())
	if err != nil {
		return nil, err
	}

	if len(overrides) > 0 {
		trainingMenuIDs := make([]training.TrainingMenuID, 0, len(overrides))
		for _, override := range overrides {
			trainingMenuIDs = append(trainingMenuIDs, override.TrainingMenuID())
		}

		trainingMenus, err := u.trainingMenuRepo.FindByIDs(ctx, trainingMenuIDs)
		if err != nil {
			return nil, err
		}

		return judgment.NewPrescriptionFromOverrides(overrides, trainingMenus), nil
	}

	trainingMenus, err := u.trainingMenuRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	elementMenus, err := u.prescriptionRuleRepo.ListElementMenus(ctx)
	if err != nil {
		return nil, err
	}

	fixedMenus, err := u.prescriptionRuleRepo.ListFixedMenus(ctx)
	if err != nil {
		return nil, err
	}

	ageDecadeMenus, err := u.prescriptionRuleRepo.ListAgeDecadeMenus(ctx)
	if err != nil {
		return nil, err
	}

	age := foundMeasurement.AgeAtMeasurement().Int()
	if _, ok := judgment.AgeDecade(age); !ok {
		log.Printf("GetJudgment: measurement %s: no age decade menus exist for age %d, prescribing without them", foundMeasurement.ID(), age)
	}

	return judgment.NewPrescription(evaluation, age, elementMenus, fixedMenus, ageDecadeMenus, trainingMenus), nil
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

	prescription, err := u.prescribe(ctx, foundMeasurement, evaluation)
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
		Prescription:  prescription,
	}
	if foundJudgment != nil {
		result.Advice = foundJudgment.Advice()
	}

	return result, nil
}

func (u *judgmentUsecase) ListOrganizationJudgments(ctx context.Context, organizationID organization.OrganizationID, includeInactive bool) ([]OrganizationJudgmentResult, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundOrganization, err := u.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, foundOrganization.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, organization.ErrOrganizationNotFound
		}
		return nil, err
	}

	customers, err := u.customerRepo.ListByOrganizationID(ctx, organizationID, includeInactive)
	if err != nil {
		return nil, err
	}

	customerIDs := make([]customer.CustomerID, 0, len(customers))
	customerGenderByID := make(map[customer.CustomerID]customer.Gender, len(customers))
	for _, foundCustomer := range customers {
		customerIDs = append(customerIDs, foundCustomer.ID())
		customerGenderByID[foundCustomer.ID()] = foundCustomer.Gender()
	}

	measurements, err := u.measurementRepo.ListByCustomerIDs(ctx, customerIDs)
	if err != nil {
		return nil, err
	}
	if len(measurements) == 0 {
		return []OrganizationJudgmentResult{}, nil
	}

	measurementItems, err := u.measurementItemRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	ageGroupStandards, err := u.ageGroupStandardRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	rankStandards, err := u.rankStandardRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	ageGroupStandardsByGender := make(map[standard.Gender][]standard.AgeGroupStandard, len(ageGroupStandards))
	for _, ageGroupStandard := range ageGroupStandards {
		gender := ageGroupStandard.Gender()
		ageGroupStandardsByGender[gender] = append(ageGroupStandardsByGender[gender], ageGroupStandard)
	}

	measurementItemByID := make(map[measurementitem.MeasurementItemID]measurementitem.MeasurementItem, len(measurementItems))
	for _, measurementItem := range measurementItems {
		measurementItemByID[measurementItem.ID()] = measurementItem
	}

	results := make([]OrganizationJudgmentResult, 0, len(measurements))
	for _, foundMeasurement := range measurements {
		age := foundMeasurement.AgeAtMeasurement().Int()
		customerGender := customerGenderByID[foundMeasurement.CustomerID()]

		gender, err := standard.NewGender(customerGender.String())
		if err != nil {
			log.Printf("ListOrganizationJudgments: measurement %s: no standards exist for customer gender %q, returning an empty evaluation", foundMeasurement.ID(), customerGender)
		}

		genderStandards := ageGroupStandardsByGender[gender]
		if err == nil && len(genderStandards) == 0 {
			log.Printf("ListOrganizationJudgments: measurement %s: no age group standards are registered for gender %q, returning an empty evaluation", foundMeasurement.ID(), gender)
		}

		measuredItems := make([]measurementitem.MeasurementItem, 0, len(foundMeasurement.Entries()))
		for _, entry := range foundMeasurement.Entries() {
			if measurementItem, ok := measurementItemByID[entry.MeasurementItemID()]; ok {
				measuredItems = append(measuredItems, measurementItem)
			}
		}

		if err == nil && len(genderStandards) > 0 {
			if unregisteredCodes := itemCodesWithoutAgeGroupStandards(measuredItems, genderStandards); len(unregisteredCodes) > 0 {
				log.Printf("ListOrganizationJudgments: measurement %s: the age group standards of gender %q are missing for items %v", foundMeasurement.ID(), gender, unregisteredCodes)
			}
		}

		if outsideCodes := itemCodesOutsideAgeGroups(measuredItems, genderStandards, age); len(outsideCodes) > 0 {
			log.Printf("ListOrganizationJudgments: measurement %s: age %d is outside the registered age groups of gender %q for items %v", foundMeasurement.ID(), age, gender, outsideCodes)
		}

		evaluation := judgment.NewEvaluation(foundMeasurement, measurementItems, gender, age, genderStandards, rankStandards)

		results = append(results, OrganizationJudgmentResult{
			CustomerID:       foundMeasurement.CustomerID(),
			MeasurementID:    foundMeasurement.ID(),
			MeasuredOn:       foundMeasurement.MeasuredOn(),
			AgeAtMeasurement: foundMeasurement.AgeAtMeasurement(),
			IsDraft:          foundMeasurement.IsDraft(),
			Evaluation:       evaluation,
		})
	}

	return results, nil
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

func (u *judgmentUsecase) findPrescribedTrainingMenus(ctx context.Context, menus []PrescribedMenuInput) ([]training.TrainingMenu, error) {
	trainingMenuIDs := make([]training.TrainingMenuID, 0, len(menus))
	seen := make(map[training.TrainingMenuID]struct{}, len(menus))
	for _, menu := range menus {
		if (menu.Element == nil) != (menu.Part == nil) {
			return nil, judgment.ErrInvalidPrescribedMenuLabels
		}
		if _, ok := seen[menu.TrainingMenuID]; ok {
			continue
		}
		seen[menu.TrainingMenuID] = struct{}{}
		trainingMenuIDs = append(trainingMenuIDs, menu.TrainingMenuID)
	}

	trainingMenus, err := u.trainingMenuRepo.FindByIDs(ctx, trainingMenuIDs)
	if err != nil {
		return nil, err
	}

	if len(trainingMenus) != len(trainingMenuIDs) {
		return nil, training.ErrTrainingMenuNotFound
	}

	trainingMenuByID := make(map[training.TrainingMenuID]training.TrainingMenu, len(trainingMenus))
	for _, trainingMenu := range trainingMenus {
		trainingMenuByID[trainingMenu.ID()] = trainingMenu
	}

	for _, menu := range menus {
		if menu.Element == nil {
			continue
		}
		trainingMenu := trainingMenuByID[menu.TrainingMenuID]
		if *menu.Element != trainingMenu.Element() || *menu.Part != trainingMenu.Part() {
			return nil, judgment.ErrInvalidPrescribedMenuLabels
		}
	}

	return trainingMenus, nil
}

func (u *judgmentUsecase) UpsertPrescription(ctx context.Context, measurementID measurement.MeasurementID, menus []PrescribedMenuInput) (judgment.Prescription, error) {
	if len(menus) == 0 {
		return nil, judgment.ErrEmptyPrescription
	}

	if _, _, err := u.findOwnedMeasurement(ctx, measurementID); err != nil {
		return nil, err
	}

	trainingMenus, err := u.findPrescribedTrainingMenus(ctx, menus)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	overrides := make([]judgment.PrescribedMenuOverride, 0, len(menus))
	for i, menu := range menus {
		sortOrder, err := training.NewSortOrder(i + 1)
		if err != nil {
			return nil, err
		}
		overrides = append(overrides, judgment.NewPrescribedMenuOverride(
			judgment.NewPrescribedMenuOverrideID(),
			measurementID,
			sortOrder,
			menu.Element,
			menu.Part,
			menu.TrainingMenuID,
			menu.Amount,
			menu.Unit,
			menu.Sets,
			now,
			now,
		))
	}

	if err := u.prescribedMenuOverrideRepo.ReplaceByMeasurementID(ctx, measurementID, overrides); err != nil {
		return nil, err
	}

	return judgment.NewPrescriptionFromOverrides(overrides, trainingMenus), nil
}

func (u *judgmentUsecase) DeletePrescription(ctx context.Context, measurementID measurement.MeasurementID) error {
	if _, _, err := u.findOwnedMeasurement(ctx, measurementID); err != nil {
		return err
	}

	return u.prescribedMenuOverrideRepo.DeleteByMeasurementID(ctx, measurementID)
}

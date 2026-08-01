package measurement

import (
	"context"
	"errors"
	"time"

	"github.com/qkitzero/fitness-service/internal/application/auth"
	"github.com/qkitzero/fitness-service/internal/application/user"
	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/measurementitem"
	"github.com/qkitzero/fitness-service/internal/domain/staff"
	"github.com/qkitzero/fitness-service/internal/domain/tenant"
)

type MeasurementEntryInput struct {
	MeasurementItemID measurementitem.MeasurementItemID
	Unmeasurable      bool
	Note              *measurement.Note
	Values            []measurement.MeasurementValue
}

type MeasurementUsecase interface {
	CreateMeasurement(ctx context.Context, customerID customer.CustomerID, measuredOn measurement.MeasuredOn, measuredBy staff.StaffID, isDraft bool, entryInputs []MeasurementEntryInput) (measurement.Measurement, error)
	GetMeasurement(ctx context.Context, measurementID measurement.MeasurementID) (measurement.Measurement, error)
	ListMeasurements(ctx context.Context, customerID customer.CustomerID) ([]measurement.Measurement, error)
	UpdateMeasurement(ctx context.Context, measurementID measurement.MeasurementID, measuredOn measurement.MeasuredOn, measuredBy staff.StaffID, isDraft bool, entryInputs []MeasurementEntryInput) (measurement.Measurement, error)
	DeleteMeasurement(ctx context.Context, measurementID measurement.MeasurementID) error
}

type measurementUsecase struct {
	authService         auth.AuthService
	userService         user.UserService
	measurementRepo     measurement.MeasurementRepository
	customerRepo        customer.CustomerRepository
	measurementItemRepo measurementitem.MeasurementItemRepository
}

func NewMeasurementUsecase(authService auth.AuthService, userService user.UserService, measurementRepo measurement.MeasurementRepository, customerRepo customer.CustomerRepository, measurementItemRepo measurementitem.MeasurementItemRepository) MeasurementUsecase {
	return &measurementUsecase{authService: authService, userService: userService, measurementRepo: measurementRepo, customerRepo: customerRepo, measurementItemRepo: measurementItemRepo}
}

func (u *measurementUsecase) verifyStaffID(ctx context.Context) (staff.StaffID, error) {
	userID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return staff.StaffID(""), err
	}

	return staff.NewStaffID(userID)
}

func (u *measurementUsecase) verifyTenantMembership(ctx context.Context, tenantID tenant.TenantID) error {
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

func (u *measurementUsecase) findOwnedCustomer(ctx context.Context, customerID customer.CustomerID) (customer.Customer, error) {
	foundCustomer, err := u.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if err := u.verifyTenantMembership(ctx, foundCustomer.TenantID()); err != nil {
		if errors.Is(err, tenant.ErrNotMember) {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, err
	}

	return foundCustomer, nil
}

func (u *measurementUsecase) findOwnedMeasurement(ctx context.Context, measurementID measurement.MeasurementID) (measurement.Measurement, customer.Customer, error) {
	foundMeasurement, err := u.measurementRepo.FindByID(ctx, measurementID)
	if err != nil {
		return nil, nil, err
	}

	foundCustomer, err := u.findOwnedCustomer(ctx, foundMeasurement.CustomerID())
	if err != nil {
		if errors.Is(err, customer.ErrCustomerNotFound) {
			return nil, nil, measurement.ErrMeasurementNotFound
		}
		return nil, nil, err
	}

	return foundMeasurement, foundCustomer, nil
}

func (u *measurementUsecase) buildEntries(ctx context.Context, entryInputs []MeasurementEntryInput) ([]measurement.MeasurementEntry, error) {
	if len(entryInputs) == 0 {
		return []measurement.MeasurementEntry{}, nil
	}

	measurementItemIDs := make([]measurementitem.MeasurementItemID, 0, len(entryInputs))
	for _, entryInput := range entryInputs {
		measurementItemIDs = append(measurementItemIDs, entryInput.MeasurementItemID)
	}

	measurementItems, err := u.measurementItemRepo.FindByIDs(ctx, measurementItemIDs)
	if err != nil {
		return nil, err
	}

	measurementItemByID := make(map[measurementitem.MeasurementItemID]measurementitem.MeasurementItem, len(measurementItems))
	for _, measurementItem := range measurementItems {
		measurementItemByID[measurementItem.ID()] = measurementItem
	}

	entries := make([]measurement.MeasurementEntry, 0, len(entryInputs))
	for _, entryInput := range entryInputs {
		measurementItem, ok := measurementItemByID[entryInput.MeasurementItemID]
		if !ok {
			return nil, measurementitem.ErrMeasurementItemNotFound
		}

		entry, err := measurement.NewMeasurementEntry(measurementItem, entryInput.Unmeasurable, entryInput.Note, entryInput.Values)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (u *measurementUsecase) CreateMeasurement(ctx context.Context, customerID customer.CustomerID, measuredOn measurement.MeasuredOn, measuredBy staff.StaffID, isDraft bool, entryInputs []MeasurementEntryInput) (measurement.Measurement, error) {
	updatedBy, err := u.verifyStaffID(ctx)
	if err != nil {
		return nil, err
	}

	foundCustomer, err := u.findOwnedCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	birthDate := foundCustomer.BirthDate()

	age, err := measurement.NewAgeAtMeasurement(birthDate.AgeAt(measuredOn.Time))
	if err != nil {
		return nil, err
	}

	entries, err := u.buildEntries(ctx, entryInputs)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	newMeasurement, err := measurement.NewMeasurement(measurement.NewMeasurementID(), customerID, measuredOn, measuredBy, age, updatedBy, isDraft, entries, now, now)
	if err != nil {
		return nil, err
	}

	if err := u.measurementRepo.Create(ctx, newMeasurement); err != nil {
		return nil, err
	}

	return newMeasurement, nil
}

func (u *measurementUsecase) GetMeasurement(ctx context.Context, measurementID measurement.MeasurementID) (measurement.Measurement, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	foundMeasurement, _, err := u.findOwnedMeasurement(ctx, measurementID)
	if err != nil {
		return nil, err
	}

	return foundMeasurement, nil
}

func (u *measurementUsecase) ListMeasurements(ctx context.Context, customerID customer.CustomerID) ([]measurement.Measurement, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	if _, err := u.findOwnedCustomer(ctx, customerID); err != nil {
		return nil, err
	}

	measurements, err := u.measurementRepo.ListByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return measurements, nil
}

func (u *measurementUsecase) UpdateMeasurement(ctx context.Context, measurementID measurement.MeasurementID, measuredOn measurement.MeasuredOn, measuredBy staff.StaffID, isDraft bool, entryInputs []MeasurementEntryInput) (measurement.Measurement, error) {
	updatedBy, err := u.verifyStaffID(ctx)
	if err != nil {
		return nil, err
	}

	foundMeasurement, foundCustomer, err := u.findOwnedMeasurement(ctx, measurementID)
	if err != nil {
		return nil, err
	}

	age := foundMeasurement.AgeAtMeasurement()
	if !foundMeasurement.MeasuredOn().Equal(measuredOn.Time) {
		birthDate := foundCustomer.BirthDate()
		if age, err = measurement.NewAgeAtMeasurement(birthDate.AgeAt(measuredOn.Time)); err != nil {
			return nil, err
		}
	}

	entries, err := u.buildEntries(ctx, entryInputs)
	if err != nil {
		return nil, err
	}

	if err := foundMeasurement.Update(measuredOn, measuredBy, age, updatedBy, isDraft, entries); err != nil {
		return nil, err
	}

	if err := u.measurementRepo.Update(ctx, foundMeasurement); err != nil {
		return nil, err
	}

	return foundMeasurement, nil
}

func (u *measurementUsecase) DeleteMeasurement(ctx context.Context, measurementID measurement.MeasurementID) error {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return err
	}

	if _, _, err := u.findOwnedMeasurement(ctx, measurementID); err != nil {
		return err
	}

	if err := u.measurementRepo.Delete(ctx, measurementID); err != nil {
		return err
	}

	return nil
}

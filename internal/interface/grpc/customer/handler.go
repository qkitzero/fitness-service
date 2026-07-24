package customer

import (
	"context"
	"errors"
	"log"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	appcustomer "github.com/qkitzero/fitness-service/internal/application/customer"
	appuser "github.com/qkitzero/fitness-service/internal/application/user"
	domaincustomer "github.com/qkitzero/fitness-service/internal/domain/customer"
)

type CustomerHandler struct {
	customerv1.UnimplementedCustomerServiceServer
	customerUsecase appcustomer.CustomerUsecase
}

func NewCustomerHandler(
	customerUsecase appcustomer.CustomerUsecase,
) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

func toDomainGender(g customerv1.Gender) (domaincustomer.Gender, error) {
	switch g {
	case customerv1.Gender_GENDER_MALE:
		return domaincustomer.GenderMale, nil
	case customerv1.Gender_GENDER_FEMALE:
		return domaincustomer.GenderFemale, nil
	case customerv1.Gender_GENDER_OTHER:
		return domaincustomer.GenderOther, nil
	default:
		return domaincustomer.NewGender("")
	}
}

func toProtoGender(g domaincustomer.Gender) customerv1.Gender {
	switch g {
	case domaincustomer.GenderMale:
		return customerv1.Gender_GENDER_MALE
	case domaincustomer.GenderFemale:
		return customerv1.Gender_GENDER_FEMALE
	case domaincustomer.GenderOther:
		return customerv1.Gender_GENDER_OTHER
	default:
		return customerv1.Gender_GENDER_UNSPECIFIED
	}
}

func toProtoBirthDate(b domaincustomer.BirthDate) *date.Date {
	return &date.Date{
		Year:  int32(b.Year()),
		Month: int32(b.Month()),
		Day:   int32(b.Day()),
	}
}

func (h *CustomerHandler) CreateCustomer(ctx context.Context, req *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	name, err := domaincustomer.NewName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	groupID, err := domaincustomer.NewGroupID(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	nameKana, err := domaincustomer.NewNameKana(req.GetNameKana())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	gender, err := toDomainGender(req.GetGender())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	birthDate, err := domaincustomer.NewBirthDate(req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.CreateCustomer(ctx, groupID, name, nameKana, gender, birthDate)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, appuser.ErrNotGroupMember) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("CreateCustomer: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &customerv1.CreateCustomerResponse{
		CustomerId: customer.ID().String(),
	}, nil
}

func (h *CustomerHandler) GetCustomer(ctx context.Context, req *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.GetCustomer(ctx, customerID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, appuser.ErrNotGroupMember) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, domaincustomer.ErrCustomerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("GetCustomer: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &customerv1.GetCustomerResponse{
		CustomerId: customer.ID().String(),
		Name:       customer.Name().String(),
		GroupId:    customer.GroupID().String(),
		NameKana:   customer.NameKana().String(),
		Gender:     toProtoGender(customer.Gender()),
		BirthDate:  toProtoBirthDate(customer.BirthDate()),
	}, nil
}

func (h *CustomerHandler) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	groupID, err := domaincustomer.NewGroupID(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customers, err := h.customerUsecase.ListCustomers(ctx, groupID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, appuser.ErrNotGroupMember) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("ListCustomers: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	customerMessages := make([]*customerv1.Customer, 0, len(customers))
	for _, c := range customers {
		customerMessages = append(customerMessages, &customerv1.Customer{
			CustomerId: c.ID().String(),
			Name:       c.Name().String(),
			GroupId:    c.GroupID().String(),
			NameKana:   c.NameKana().String(),
			Gender:     toProtoGender(c.Gender()),
			BirthDate:  toProtoBirthDate(c.BirthDate()),
		})
	}

	return &customerv1.ListCustomersResponse{
		Customers: customerMessages,
	}, nil
}

func (h *CustomerHandler) UpdateCustomer(ctx context.Context, req *customerv1.UpdateCustomerRequest) (*customerv1.UpdateCustomerResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := domaincustomer.NewName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	nameKana, err := domaincustomer.NewNameKana(req.GetNameKana())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	gender, err := toDomainGender(req.GetGender())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	birthDate, err := domaincustomer.NewBirthDate(req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.UpdateCustomer(ctx, customerID, name, nameKana, gender, birthDate)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, appuser.ErrNotGroupMember) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, domaincustomer.ErrCustomerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("UpdateCustomer: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &customerv1.UpdateCustomerResponse{
		CustomerId: customer.ID().String(),
		Name:       customer.Name().String(),
		GroupId:    customer.GroupID().String(),
		NameKana:   customer.NameKana().String(),
		Gender:     toProtoGender(customer.Gender()),
		BirthDate:  toProtoBirthDate(customer.BirthDate()),
	}, nil
}

func (h *CustomerHandler) DeleteCustomer(ctx context.Context, req *customerv1.DeleteCustomerRequest) (*customerv1.DeleteCustomerResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.customerUsecase.DeleteCustomer(ctx, customerID); err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, appuser.ErrNotGroupMember) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, domaincustomer.ErrCustomerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("DeleteCustomer: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &customerv1.DeleteCustomerResponse{}, nil
}

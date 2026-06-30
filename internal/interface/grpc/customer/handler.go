package customer

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	appcustomer "github.com/qkitzero/fitness-service/internal/application/customer"
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

func (h *CustomerHandler) CreateCustomer(ctx context.Context, req *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	name, err := domaincustomer.NewName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.CreateCustomer(ctx, name)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
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
		if errors.Is(err, domaincustomer.ErrCustomerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("GetCustomer: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &customerv1.GetCustomerResponse{
		CustomerId: customer.ID().String(),
		Name:       customer.Name().String(),
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

	customer, err := h.customerUsecase.UpdateCustomer(ctx, customerID, name)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
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
		if errors.Is(err, domaincustomer.ErrCustomerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("DeleteCustomer: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &customerv1.DeleteCustomerResponse{}, nil
}

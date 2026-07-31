package customer

import (
	"context"
	"errors"
	"log"
	"strings"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	appcustomer "github.com/qkitzero/fitness-service/internal/application/customer"
	domainaddress "github.com/qkitzero/fitness-service/internal/domain/address"
	domaincontact "github.com/qkitzero/fitness-service/internal/domain/contact"
	domaincustomer "github.com/qkitzero/fitness-service/internal/domain/customer"
	domainorganization "github.com/qkitzero/fitness-service/internal/domain/organization"
	domaintenant "github.com/qkitzero/fitness-service/internal/domain/tenant"
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

func toProtoCustomer(c domaincustomer.Customer) *customerv1.Customer {
	msg := &customerv1.Customer{
		CustomerId: c.ID().String(),
		Name:       c.Name().String(),
		TenantId:   c.TenantID().String(),
		NameKana:   c.NameKana().String(),
		Gender:     toProtoGender(c.Gender()),
		BirthDate:  toProtoBirthDate(c.BirthDate()),
		IsActive:   c.IsActive(),
	}
	if v := c.Phone(); v != nil {
		s := v.String()
		msg.Phone = &s
	}
	if v := c.Email(); v != nil {
		s := v.String()
		msg.Email = &s
	}
	addr := c.Address()
	if v := addr.PostalCode(); v != nil {
		s := v.String()
		msg.PostalCode = &s
	}
	if v := addr.Prefecture(); v != nil {
		s := v.String()
		msg.Prefecture = &s
	}
	if v := addr.City(); v != nil {
		s := v.String()
		msg.City = &s
	}
	if v := addr.Street(); v != nil {
		s := v.String()
		msg.Street = &s
	}
	if v := addr.Building(); v != nil {
		s := v.String()
		msg.Building = &s
	}
	if v := c.EmergencyContactName(); v != nil {
		s := v.String()
		msg.EmergencyContactName = &s
	}
	if v := c.EmergencyContactRelationship(); v != nil {
		s := v.String()
		msg.EmergencyContactRelationship = &s
	}
	if v := c.EmergencyContactPhone(); v != nil {
		s := v.String()
		msg.EmergencyContactPhone = &s
	}
	if v := c.OrganizationID(); v != nil {
		s := v.String()
		msg.OrganizationId = &s
	}
	return msg
}

type customerFieldsRequest interface {
	GetName() string
	GetNameKana() string
	GetGender() customerv1.Gender
	GetBirthDate() *date.Date
	GetPhone() string
	GetEmail() string
	GetPostalCode() string
	GetPrefecture() string
	GetCity() string
	GetStreet() string
	GetBuilding() string
	GetEmergencyContactName() string
	GetEmergencyContactRelationship() string
	GetEmergencyContactPhone() string
	GetOrganizationId() string
}

type customerFields struct {
	name                         domaincustomer.Name
	nameKana                     domaincustomer.NameKana
	gender                       domaincustomer.Gender
	birthDate                    domaincustomer.BirthDate
	phone                        *domaincontact.Phone
	email                        *domaincontact.Email
	address                      domainaddress.Address
	emergencyContactName         *domaincustomer.EmergencyContactName
	emergencyContactRelationship *domaincustomer.EmergencyContactRelationship
	emergencyContactPhone        *domaincontact.Phone
	organizationID               *domainorganization.OrganizationID
}

func parseAddress(req customerFieldsRequest) (domainaddress.Address, error) {
	postalCode, err := domainaddress.NewPostalCode(req.GetPostalCode())
	if err != nil {
		return domainaddress.Address{}, err
	}
	prefecture, err := domainaddress.NewPrefecture(req.GetPrefecture())
	if err != nil {
		return domainaddress.Address{}, err
	}
	city, err := domainaddress.NewCity(req.GetCity())
	if err != nil {
		return domainaddress.Address{}, err
	}
	street, err := domainaddress.NewStreet(req.GetStreet())
	if err != nil {
		return domainaddress.Address{}, err
	}
	building, err := domainaddress.NewBuilding(req.GetBuilding())
	if err != nil {
		return domainaddress.Address{}, err
	}
	return domainaddress.NewAddress(postalCode, prefecture, city, street, building), nil
}

func parseCustomerFields(req customerFieldsRequest) (customerFields, error) {
	var f customerFields
	var err error
	if f.name, err = domaincustomer.NewName(req.GetName()); err != nil {
		return f, err
	}
	if f.nameKana, err = domaincustomer.NewNameKana(req.GetNameKana()); err != nil {
		return f, err
	}
	if f.gender, err = toDomainGender(req.GetGender()); err != nil {
		return f, err
	}
	if f.birthDate, err = domaincustomer.NewBirthDate(req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay()); err != nil {
		return f, err
	}
	if f.phone, err = domaincontact.NewPhone(req.GetPhone()); err != nil {
		return f, err
	}
	if f.email, err = domaincontact.NewEmail(req.GetEmail()); err != nil {
		return f, err
	}
	if f.address, err = parseAddress(req); err != nil {
		return f, err
	}
	if f.emergencyContactName, err = domaincustomer.NewEmergencyContactName(req.GetEmergencyContactName()); err != nil {
		return f, err
	}
	if f.emergencyContactRelationship, err = domaincustomer.NewEmergencyContactRelationship(req.GetEmergencyContactRelationship()); err != nil {
		return f, err
	}
	if f.emergencyContactPhone, err = domaincontact.NewPhone(req.GetEmergencyContactPhone()); err != nil {
		return f, err
	}
	if s := strings.TrimSpace(req.GetOrganizationId()); s != "" {
		organizationID, err := domainorganization.NewOrganizationIDFromString(s)
		if err != nil {
			return f, err
		}
		f.organizationID = &organizationID
	}
	return f, nil
}

func mapCustomerError(err error, op string) error {
	if errors.Is(err, domaintenant.ErrNotMember) {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	if errors.Is(err, domaincustomer.ErrCustomerNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, domaincustomer.ErrOrganizationNotInTenant) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.Unauthenticated, codes.PermissionDenied:
			return err
		}
	}
	log.Printf("%s: internal error: %v", op, err)
	return status.Error(codes.Internal, "internal error")
}

func (h *CustomerHandler) CreateCustomer(ctx context.Context, req *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	tenantID, err := domaintenant.NewTenantID(req.GetTenantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fields, err := parseCustomerFields(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.CreateCustomer(ctx, tenantID, fields.name, fields.nameKana, fields.gender, fields.birthDate, fields.phone, fields.email, fields.address, fields.emergencyContactName, fields.emergencyContactRelationship, fields.emergencyContactPhone, fields.organizationID)
	if err != nil {
		return nil, mapCustomerError(err, "CreateCustomer")
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
		return nil, mapCustomerError(err, "GetCustomer")
	}

	return &customerv1.GetCustomerResponse{
		Customer: toProtoCustomer(customer),
	}, nil
}

func (h *CustomerHandler) ListCustomers(ctx context.Context, req *customerv1.ListCustomersRequest) (*customerv1.ListCustomersResponse, error) {
	tenantID, err := domaintenant.NewTenantID(req.GetTenantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customers, err := h.customerUsecase.ListCustomers(ctx, tenantID, req.GetIncludeInactive())
	if err != nil {
		return nil, mapCustomerError(err, "ListCustomers")
	}

	customerMessages := make([]*customerv1.Customer, 0, len(customers))
	for _, c := range customers {
		customerMessages = append(customerMessages, toProtoCustomer(c))
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
	fields, err := parseCustomerFields(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.UpdateCustomer(ctx, customerID, fields.name, fields.nameKana, fields.gender, fields.birthDate, fields.phone, fields.email, fields.address, fields.emergencyContactName, fields.emergencyContactRelationship, fields.emergencyContactPhone, fields.organizationID)
	if err != nil {
		return nil, mapCustomerError(err, "UpdateCustomer")
	}

	return &customerv1.UpdateCustomerResponse{
		Customer: toProtoCustomer(customer),
	}, nil
}

func (h *CustomerHandler) SetCustomerActive(ctx context.Context, req *customerv1.SetCustomerActiveRequest) (*customerv1.SetCustomerActiveResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := h.customerUsecase.SetCustomerActive(ctx, customerID, req.GetIsActive())
	if err != nil {
		return nil, mapCustomerError(err, "SetCustomerActive")
	}

	return &customerv1.SetCustomerActiveResponse{
		Customer: toProtoCustomer(customer),
	}, nil
}

func (h *CustomerHandler) DeleteCustomer(ctx context.Context, req *customerv1.DeleteCustomerRequest) (*customerv1.DeleteCustomerResponse, error) {
	customerID, err := domaincustomer.NewCustomerIDFromString(req.GetCustomerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.customerUsecase.DeleteCustomer(ctx, customerID); err != nil {
		return nil, mapCustomerError(err, "DeleteCustomer")
	}

	return &customerv1.DeleteCustomerResponse{}, nil
}

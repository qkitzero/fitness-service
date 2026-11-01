package judgment

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/customer"
	"github.com/qkitzero/fitness-service/internal/domain/measurement"
	"github.com/qkitzero/fitness-service/internal/domain/organization"
)

const (
	DefaultOrganizationJudgmentPageSize = 500
	MaxOrganizationJudgmentPageSize     = 1000
)

var (
	ErrInvalidPageSize  = errors.New("invalid page size")
	ErrInvalidPageToken = errors.New("invalid page token")
)

type OrganizationJudgmentPage struct {
	Judgments     []OrganizationJudgmentResult
	NextPageToken string
}

type organizationJudgmentPageToken struct {
	OrganizationID  string `json:"o"`
	IncludeInactive bool   `json:"i"`
	CustomerID      string `json:"c"`
	MeasuredOn      string `json:"d"`
	MeasurementID   string `json:"m"`
}

func resolvePageSize(pageSize int) (int, error) {
	switch {
	case pageSize < 0:
		return 0, ErrInvalidPageSize
	case pageSize == 0:
		return DefaultOrganizationJudgmentPageSize, nil
	case pageSize > MaxOrganizationJudgmentPageSize:
		return MaxOrganizationJudgmentPageSize, nil
	default:
		return pageSize, nil
	}
}

func encodePageToken(organizationID organization.OrganizationID, includeInactive bool, last measurement.Measurement) (string, error) {
	data, err := json.Marshal(organizationJudgmentPageToken{
		OrganizationID:  organizationID.String(),
		IncludeInactive: includeInactive,
		CustomerID:      last.CustomerID().String(),
		MeasuredOn:      last.MeasuredOn().Format(time.DateOnly),
		MeasurementID:   last.ID().String(),
	})
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodePageToken(pageToken string, organizationID organization.OrganizationID, includeInactive bool) (*measurement.ListCursor, error) {
	if pageToken == "" {
		return nil, nil
	}

	data, err := base64.RawURLEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, ErrInvalidPageToken
	}

	var token organizationJudgmentPageToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, ErrInvalidPageToken
	}

	if token.OrganizationID != organizationID.String() || token.IncludeInactive != includeInactive {
		return nil, ErrInvalidPageToken
	}

	customerID, err := customer.NewCustomerIDFromString(token.CustomerID)
	if err != nil {
		return nil, ErrInvalidPageToken
	}

	measuredOn, err := time.Parse(time.DateOnly, token.MeasuredOn)
	if err != nil {
		return nil, ErrInvalidPageToken
	}

	measurementID, err := measurement.NewMeasurementIDFromString(token.MeasurementID)
	if err != nil {
		return nil, ErrInvalidPageToken
	}

	return &measurement.ListCursor{
		CustomerID:    customerID,
		MeasuredOn:    measurement.MeasuredOn{Time: measuredOn},
		MeasurementID: measurementID,
	}, nil
}

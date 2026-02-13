package store

import (
	"github.com/google/uuid"
	tenantdom "github.com/iamonah/merchcore/internal/domain/tenant"
)

type CreateTenantRequest struct {
	BusinessName      string        `json:"business_name" validate:"required"`
	Description       string        `json:"description" validate:"required"`
	Subdomain         *string       `json:"subdomain,omitempty"`
	BusinessMode      string        `json:"business_mode" validate:"required"`
	BusinessCategory  string        `json:"business_category" validate:"required"`
	Plan              string        `json:"plan" validate:"required"`
	Domain            *string       `json:"domain,omitempty"`
	LogoURL           *string       `json:"logo_url,omitempty"`
	BusinessAddress   *AddressInput `json:"business_address,omitempty"`
	BillingAddress    *AddressInput `json:"billing_address,omitempty"`
	NumberOfEmployees int32         `json:"number_of_employees,omitempty"`
}

type AddressInput struct {
	Street     string `json:"street" validate:"required"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
	IsDefault  bool   `json:"is_default,omitempty"`
	Type       string `json:"type" validate:"required"`
}

func toAddress(userID uuid.UUID, in *AddressInput) *tenantdom.AddressInput {
	if in == nil {
		return nil
	}
	return &tenantdom.AddressInput{
		UserID:     userID,
		Street:     in.Street,
		City:       in.City,
		State:      in.State,
		PostalCode: in.PostalCode,
		Country:    in.Country,
		IsDefault:  in.IsDefault,
		Type:       in.Type,
	}
}

func toTenantCreateInput(userID uuid.UUID, req CreateTenantRequest) tenantdom.CreateTenant {
	return tenantdom.CreateTenant{
		UserID:            userID,
		BusinessName:      req.BusinessName,
		Description:       req.Description,
		Subdomain:         req.Subdomain,
		Domain:            req.Domain,
		LogoURL:           req.LogoURL,
		BusinessMode:      req.BusinessMode,
		BusinessCategory:  req.BusinessCategory,
		Plan:              req.Plan,
		BusinessAddress:   toAddress(userID, req.BusinessAddress),
		BillingAddress:    toAddress(userID, req.BillingAddress),
		NumberOfEmployees: req.NumberOfEmployees,
	}
}

type CreateTenantResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	BusinessName      string    `json:"business_name"`
	Domain            string    `json:"domain"`
	Subdomain         string    `json:"subdomain"`
	Plan              string    `json:"plan"`
	Status            string    `json:"status"`
	BusinessMode      string    `json:"business_mode"`
	NumberOfEmployees int32     `json:"number_of_employees"`
}

func toCreateTenantResponse(t *tenantdom.TenantProfile) CreateTenantResponse {
	return CreateTenantResponse{
		ID:                t.ID,
		UserID:            t.UserID,
		BusinessName:      t.BusinessName,
		Domain:            *t.Domain,
		Subdomain:         *t.Subdomain,
		Plan:              string(t.Plan),
		Status:            string(t.Status),
		BusinessMode:      string(t.BusinessMode),
		NumberOfEmployees: t.NumberOfEmployees,
	}
}

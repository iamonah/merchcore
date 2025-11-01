package tenant

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/iamonah/merchcore/internal/domain/shared/address"
	"github.com/iamonah/merchcore/internal/sdk/errs"

	"github.com/google/uuid"
)

type TenantProfile struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	BusinessName string
	Description  string
	Subdomain    *string
	Domain       *string
	LogoURL      string
	Status       TenantStatus
	Plan         PlanType
	BusinessMode BusinessMode
	Addresses    address.Addresses
	TrialStartAt *time.Time
	TrialEndAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// if logo is uploaded upload to s3 and get url here
// client pays for store and a store is created
func NewTenantProfile(storeInfo CreateTenant) (*TenantProfile, error) {
	fieldErrs := errs.NewFieldErrors()

	businessName := strings.TrimSpace(storeInfo.BusinessName)
	if businessName == "" {
		fieldErrs.AddFieldError("business_name", errors.New("business name required"))
	}

	description := strings.TrimSpace(storeInfo.Description)
	if description == "" {
		fieldErrs.AddFieldError("description", errors.New("description required"))
	} else if utf8.RuneCountInString(description) > 255 {
		fieldErrs.AddFieldError("description", errors.New("cannot be more than 255 characters"))
	}

	var subdomain string
	if storeInfo.Subdomain == nil {
		subdomain = strings.ToLower(strings.ReplaceAll(businessName, " ", "")) + ".merchcore.com"
	} else {
		subdomain = *storeInfo.Subdomain
	}

	plan, err := ParsePlanType(storeInfo.Plan)
	if err != nil {
		fieldErrs.AddFieldError("plan", errors.New("invalid plan"))
	}

	mode, err := ParseBusinessMode(storeInfo.BusinessMode)
	if err != nil {
		fieldErrs.AddFieldError("business_mode", errors.New("invalid business mode"))
	}

	var addresses address.Addresses

	if mode == BusinessModeHybrid {
		if storeInfo.BusinessAddress == nil {
			fieldErrs.AddFieldError("business_address", errors.New("business address required for hybrid mode"))
		} else {
			addr, err := NewAddress(*storeInfo.BusinessAddress)
			if err != nil {
				fieldErrs.AddFieldError("business_address", err)
			} else {
				addresses.AddAddresses(*addr)
			}
		}
	}

	if plan != FreePlan {
		if storeInfo.BillingAddress == nil {
			fieldErrs.AddFieldError("billing_address", errors.New("billing address required for paid plans"))
		} else {
			addr, err := NewAddress(*storeInfo.BillingAddress)
			if err != nil {
				fieldErrs.AddFieldError("billing_address", err)
			} else {
				addresses.AddAddresses(*addr)
			}
		}
	}

	if fieldErrs.ToError() != nil {
		return nil, fieldErrs
	}

	now := time.Now()
	var trialEnd *time.Time

	if plan == FreePlan {
		end := now.Add(15 * 24 * time.Hour)
		trialEnd = &end
	}

	tenant := &TenantProfile{
		ID:           uuid.New(),
		UserID:       storeInfo.UserID,
		BusinessName: businessName,
		Description:  description,
		Subdomain:    &subdomain,
		Domain:       storeInfo.Domain,
		LogoURL:      *storeInfo.LogoURL,
		Status:       TenantStatusMaintenance,
		Plan:         plan,
		BusinessMode: mode,
		Addresses:    addresses,
		TrialStartAt: &now,
		TrialEndAt:   trialEnd,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return tenant, nil
}

// // Example methods for invariants
// func (t *TenantProfile) IsTrialActive(now time.Time) bool {
// 	return t.Status == TenantStatusTrial && now.Before(*t.TrialEndAt)
// }

// func (t *TenantProfile) UpgradeToActive(domain string) {
// 	t.Status = TenantStatusActive
// 	t.Domain = domain
// 	t.UpdatedAt = time.Now()
// 	// Trigger hooks in service: provision custom domain, etc.
// }

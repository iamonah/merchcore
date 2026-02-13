package tenant

import (
	"context"
	"errors"
	"fmt"

	"github.com/iamonah/merchcore/internal/infra/database"
	"github.com/iamonah/merchcore/internal/sdk/errs"
)

type TenantBusiness struct {
	storer TenantRepository
	trx    database.TransactorTX
}

type TenantBusinessCfg func(tb *TenantBusiness) error

func NewTenantBusiness(cfgs ...TenantBusinessCfg) (*TenantBusiness, error) {
	tb := &TenantBusiness{}
	for _, cfg := range cfgs {
		if err := cfg(tb); err != nil {
			return nil, err
		}
	}
	if tb.storer == nil {
		return nil, errors.New("tenant repository is required")
	}
	if tb.trx == nil {
		return nil, errors.New("transaction manager is required")
	}
	return tb, nil
}

func WithTenantRepository(st TenantRepository) TenantBusinessCfg {
	return func(tb *TenantBusiness) error {
		tb.storer = st
		return nil
	}
}

func WithTransactor(trx database.TransactorTX) TenantBusinessCfg {
	return func(tb *TenantBusiness) error {
		tb.trx = trx
		return nil
	}
}

func (tb *TenantBusiness) CreateTenant(ctx context.Context, input CreateTenant) (*TenantProfile, error) {
	tenantProfile, err := NewTenantProfile(input)
	if err != nil {
		return nil, errs.NewDomainError(errs.InvalidArgument, err)
	}

	subdomainTaken, err := tb.storer.CheckSubdomainAvailability(ctx, *tenantProfile.Subdomain)
	if err != nil {
		return nil, fmt.Errorf("checksubdomainavailability: %w", err)
	}
	if subdomainTaken {
		return nil, errs.NewDomainError(errs.AlreadyExists, ErrSubDomain)
	}

	if err := tb.trx.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := tb.storer.CreateTenant(txCtx, tenantProfile); err != nil {
			return err
		}
		if err := tb.storer.CreateTenantSchema(txCtx, tenantProfile.UserID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		switch {
		case errors.Is(err, ErrDomain), errors.Is(err, ErrSubDomain):
			return nil, errs.NewDomainError(errs.AlreadyExists, err)
		case errors.Is(err, ErrInvalidEnumValue):
			return nil, errs.NewDomainError(errs.InvalidArgument, err)
		default:
			return nil, fmt.Errorf("tenantcreate-trx: %w", err)
		}
	}

	return tenantProfile, nil
}

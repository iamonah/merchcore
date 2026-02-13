package store

import (
	"errors"

	"github.com/iamonah/merchcore/internal/domain/tenant"
	"github.com/iamonah/merchcore/internal/sdk/jobs"
	"github.com/rs/zerolog"
)

type TenantService struct {
	log     *zerolog.Logger
	job     jobs.JobService
	tenants *tenant.TenantBusiness
}

type TenantConfiguration func(ts *TenantService) error

func NewTenantService(cfgs ...TenantConfiguration) (*TenantService, error) {
	ts := &TenantService{}
	for _, cfg := range cfgs {
		if err := cfg(ts); err != nil {
			return nil, err
		}
	}
	if ts.log == nil {
		return nil, errors.New("logger is required")
	}
	if ts.tenants == nil {
		return nil, errors.New("tenant business is required")
	}
	return ts, nil
}

// NewUserService is kept for backward compatibility with existing call sites.
func NewUserService(cfgs ...TenantConfiguration) (*TenantService, error) {
	return NewTenantService(cfgs...)
}

func WithTenantBusiness(tb *tenant.TenantBusiness) TenantConfiguration {
	return func(ts *TenantService) error {
		ts.tenants = tb
		return nil
	}
}

func WithLog(log *zerolog.Logger) TenantConfiguration {
	return func(ts *TenantService) error {
		ts.log = log
		return nil
	}
}

func WithJob(job jobs.JobService) TenantConfiguration {
	return func(ts *TenantService) error {
		ts.job = job
		return nil
	}
}

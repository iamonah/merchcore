package auth

import (
	"github.com/IamOnah/storefronthq/internal/domain/users"
	"github.com/IamOnah/storefronthq/internal/infra/database"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/jobs"

	"github.com/rs/zerolog"
)

type UserService struct {
	users users.UserRepository
	auth  authz.JWTAuthMaker
	log   *zerolog.Logger
	job   jobs.JobService
	trx   database.TransactorTX
}

type UserConfiguration func(us *UserService) error

func NewUserService(cfgs ...UserConfiguration) (*UserService, error) {
	os := &UserService{}
	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

func WithUserRepository(ur users.UserRepository) UserConfiguration {
	return func(us *UserService) error {
		us.users = ur
		return nil
	}
}

func WithAuth(auth authz.JWTAuthMaker) UserConfiguration {
	return func(us *UserService) error {
		us.auth = auth
		return nil
	}
}

func WithLog(log *zerolog.Logger) UserConfiguration {
	return func(us *UserService) error {
		us.log = log
		return nil
	}
}

func WithJob(job *jobs.JobClient) UserConfiguration {
	return func(us *UserService) error {
		us.job = job
		return nil
	}
}

func WithTrxManager(trxManager *database.TRXManager) UserConfiguration {
	return func(us *UserService) error {
		us.trx = trxManager
		return nil
	}
}

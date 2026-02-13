package auth

import (
	"github.com/iamonah/merchcore/internal/domain/users"
	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/jobs"

	"github.com/rs/zerolog"
)

type UserService struct {
	auth  authz.JWTAuthMaker
	log   *zerolog.Logger
	job   jobs.JobService
	users users.ExtUserBusiness
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

func WithUserBusiness(ub *users.UserBusiness) UserConfiguration {
	return func(us *UserService) error {
		us.users = ub
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

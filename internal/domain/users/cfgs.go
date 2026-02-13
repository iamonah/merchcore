package users

import (
	"github.com/iamonah/merchcore/internal/config"
	"github.com/iamonah/merchcore/internal/infra/cache"
	"github.com/iamonah/merchcore/internal/infra/database"
	"github.com/iamonah/merchcore/internal/sdk/authz"
)

func WithUserRepository(store UserRepository) UserBusinessCfg {
	return func(ub *UserBusiness) error {
		ub.storer = store
		return nil
	}
}

func WithTrxManager(trx *database.TRXManager) UserBusinessCfg {
	return func(ub *UserBusiness) error {
		ub.trx = trx
		return nil
	}
}

func WithAuthz(authz *authz.JWTAuthMaker) UserBusinessCfg {
	return func(ub *UserBusiness) error {
		ub.authz = authz
		return nil
	}
}

func WithConfigs(cfg *config.Config) UserBusinessCfg {
	return func(ub *UserBusiness) error {
		ub.config = cfg
		return nil
	}
}

func WithCache(cache cache.Cache) UserBusinessCfg {
	return func(ub *UserBusiness) error {
		ub.cache = cache
		return nil
	}
}

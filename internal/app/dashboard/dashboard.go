package dashboard

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/domain/users"
	"github.com/rs/zerolog"
)

type DashboardService struct {
	log   *zerolog.Logger
	users *users.UserBusiness
}

func NewDashboardService() *DashboardService {
	return nil
}

func (d *DashboardService) GetOverview(ctx context.Context, userID uuid.UUID) error {
	return nil
}

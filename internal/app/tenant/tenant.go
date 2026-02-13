package store

import (
	"fmt"
	"net/http"

	"github.com/iamonah/merchcore/internal/sdk/base"
	"github.com/iamonah/merchcore/internal/sdk/errs"
)

func (ts *TenantService) CreateTenant(w http.ResponseWriter, r *http.Request) error {
	reqID, err := base.GetReqIDCTX(r)
	if err != nil {
		return errs.Newf(errs.Internal, "getreqidCTX: %s", err)
	}
	pl, err := base.GetJWTPayloadCTX(r)
	if err != nil {
		return errs.New(errs.Unauthenticated, fmt.Errorf("unauthorized"))
	}

	var req CreateTenantRequest
	if err := base.ReadJSON(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	if err := errs.NewValidate(req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	profile, err := ts.tenants.CreateTenant(r.Context(), toTenantCreateInput(pl.UserID, req))
	if err != nil {
		if derr, ok := errs.IsDomainError(err); ok {
			return errs.New(derr.Code, derr)
		}
		return errs.Newf(errs.Internal, "createtenant: reqID[%s] userID[%s]: %s", reqID, pl.UserID, err)
	}

	ts.log.Info().
		Str("event", "tenant.create").
		Str("req_id", reqID).
		Str("tenant_id", profile.ID.String()).
		Str("user_id", pl.UserID.String()).
		Msg("tenant created")

	if err := base.WriteJSON(w, http.StatusCreated, toCreateTenantResponse(profile)); err != nil {
		return errs.Newf(errs.Internal, "writejson: %s", err)
	}
	return nil
}

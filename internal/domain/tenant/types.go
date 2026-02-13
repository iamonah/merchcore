package tenant

import (
	"fmt"
	"strings"
)

type TenantStatus string

var tenantStatuses = make(map[string]TenantStatus)

func newTenantStatus(v string) TenantStatus {
	ts := TenantStatus(v)
	tenantStatuses[strings.ToLower(v)] = ts
	return ts
}

var (
	TenantStatusActive      = newTenantStatus("active")
	TenantStatusMaintenance = newTenantStatus("maintenance")
	TenantStatusSuspended   = newTenantStatus("suspended")
	TenantStatusArchived    = newTenantStatus("archived")
)

func ParseTenantStatus(v string) (TenantStatus, error) {
	ts, ok := tenantStatuses[strings.ToLower(v)]
	if !ok {
		return "", fmt.Errorf("invalid tenant status: %v", v)
	}
	return ts, nil
}

type PlanType string

var planTypes = make(map[string]PlanType)

func newPlanType(v string) PlanType {
	pt := PlanType(v)
	planTypes[strings.ToLower(v)] = pt
	return pt
}

var (
	FreePlan       = newPlanType("free")
	CorePlan       = newPlanType("core")
	ProPlan        = newPlanType("pro")
	EnterprisePlan = newPlanType("enterprise")
)

func ParsePlanType(v string) (PlanType, error) {
	pt, ok := planTypes[strings.ToLower(v)]
	if !ok {
		return "", fmt.Errorf("invalid plan type: %v", v)
	}
	return pt, nil
}

type BusinessMode string

var businessModes = make(map[string]BusinessMode)

func newBusinessMode(v string) BusinessMode {
	bm := BusinessMode(v)
	businessModes[strings.ToLower(v)] = bm
	return bm
}

var (
	BusinessModeOnline = newBusinessMode("online")
	BusinessModeHybrid = newBusinessMode("hybrid")
)

func ParseBusinessMode(v string) (BusinessMode, error) {
	bm, ok := businessModes[strings.ToLower(v)]
	if !ok {
		return "", fmt.Errorf("invalid business mode: %v", v)
	}
	return bm, nil
}

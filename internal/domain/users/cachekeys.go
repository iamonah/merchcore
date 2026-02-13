package users

import "fmt"

func SessionAccessKeys(userID any) string {
	return fmt.Sprintf("access:%v", userID)
}
func SessionRefreshKeys(userID any) string {
	return fmt.Sprintf("refresh:%v", userID)
}

func UserProfile(userID any) string {
	return fmt.Sprintf("user:profile:%v", userID)
}

func TenantSettings(tenantID any) string {
	return fmt.Sprintf("tenant:settings:%v", tenantID)
}

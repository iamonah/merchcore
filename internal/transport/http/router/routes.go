package router

import (
	"net/http"

	"github.com/iamonah/merchcore/internal/app/auth"
	"github.com/iamonah/merchcore/internal/sdk/authz"
	"github.com/iamonah/merchcore/internal/sdk/midd"
	"github.com/rs/zerolog"
)

func SetupRouter(
	us *auth.UserService,
	log *zerolog.Logger,
	maker *authz.JWTAuthMaker,
	// te *store.TenantService,
) http.Handler {
	app := NewApp(log, midd.RecoverPanic(log))

	authbearer := midd.AuthBearer(maker)
	// version := "1"
	app.HandleFunc(http.MethodPost, "/auth/signin", us.Authenticate)
	app.HandleFunc(http.MethodPost, "/auth/register", us.RegisterUser)
	app.HandleFunc(http.MethodPost, "/auth/signout", us.SignOut, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/reset-password", us.ResetPassword)
	app.HandleFunc(http.MethodPost, "/auth/forgot-password", us.ForgotPassword)
	app.HandleFunc(http.MethodPost, "/auth/activate", us.ActivateUser, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/change-password", us.ChangePassword, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/token/renew-access", us.RenewAccessToken, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/resend-token", us.ResendVerificationToken, authbearer)

	// app.HandleFunc(http.MethodGet, "/api/stores/:id", us.GetStore)
	// app.HandleFunc(http.MethodPut, "/api/stores/:id", us.UpdateStore)
	// app.HandleFunc(http.MethodDelete, "/api/stores/:id", us.DeleteStore)
	// app.HandleFunc(http.MethodPost, "/api/create-stores/", te.CreateTenant)

	// 🏬 Store Management
	// ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/stores", ds.ListStores, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/stores", ds.CreateStore, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/stores/:id", ds.GetStore, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/stores/:id", ds.UpdateStore, authbearer)
	// app.HandleFunc(http.MethodDelete, "/dashboard/stores/:id", ds.DeleteStore, authbearer)

	// ------------------------------
	// // 🎨 Appearance / Customization
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/themes", ds.ListThemes, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/stores/:id/theme", ds.ApplyTheme, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/stores/:id/logo", ds.UploadLogo, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/stores/:id/theme/preview", ds.PreviewTheme, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/stores/:id/theme/colors", ds.UpdateThemeColors, authbearer)

	// ------------------------------
	// // 🛍️ Products & Inventory
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/products", ds.ListProducts, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/products", ds.CreateProduct, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/products/:id", ds.GetProduct, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/products/:id", ds.UpdateProduct, authbearer)
	// app.HandleFunc(http.MethodDelete, "/dashboard/products/:id", ds.DeleteProduct, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/products/:id/stock", ds.UpdateStock, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/products/import", ds.ImportProductsCSV, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/products/export", ds.ExportProductsCSV, authbearer)

	// // ------------------------------
	// // 📦 Orders & Fulfillment
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/orders", ds.ListOrders, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/orders/:id", ds.GetOrder, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/orders", ds.CreateOrderManual, authbearer) // Manual creation
	// app.HandleFunc(http.MethodPut, "/dashboard/orders/:id/status", ds.UpdateOrderStatus, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/orders/export", ds.ExportOrdersCSV, authbearer)

	// // ------------------------------
	// // 💰 Finance / Payouts
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/finance/transactions", ds.ListTransactions, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/finance/payouts", ds.ListPayouts, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/finance/withdraw", ds.RequestWithdrawal, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/finance/overview", ds.GetFinanceOverview, authbearer)

	// // ------------------------------
	// // 👥 Customers
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/customers", ds.ListCustomers, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/customers", ds.AddCustomer, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/customers/:id", ds.GetCustomer, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/customers/:id", ds.UpdateCustomer, authbearer)
	// app.HandleFunc(http.MethodDelete, "/dashboard/customers/:id", ds.DeleteCustomer, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/customers/export", ds.ExportCustomersCSV, authbearer)

	// // ------------------------------
	// // 📊 Analytics
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/analytics/sales", ds.GetSalesAnalytics, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/analytics/traffic", ds.GetTrafficAnalytics, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/analytics/products", ds.GetProductAnalytics, authbearer)

	// // ------------------------------
	// // ⚙️ Settings & Integrations
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/settings", ds.GetSettings, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/settings", ds.UpdateSettings, authbearer)
	// app.HandleFunc(http.MethodGet, "/dashboard/integrations", ds.ListIntegrations, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/integrations/:type", ds.ConnectIntegration, authbearer)
	// app.HandleFunc(http.MethodDelete, "/dashboard/integrations/:type", ds.DisconnectIntegration, authbearer)

	// // ------------------------------
	// // 🧾 Export / Import (Global)
	// // ------------------------------
	// app.HandleFunc(http.MethodGet, "/dashboard/export/:entity", ds.ExportData, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/import/:entity", ds.ImportData, authbearer)

	// // 👩‍💻 Team / Staff Management
	// app.HandleFunc(http.MethodGet, "/dashboard/team", ds.ListTeamMembers, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/team", ds.InviteTeamMember, authbearer)
	// app.HandleFunc(http.MethodPut, "/dashboard/team/:id/role", ds.UpdateTeamMemberRole, authbearer)
	// app.HandleFunc(http.MethodDelete, "/dashboard/team/:id", ds.RemoveTeamMember, authbearer)
	// app.HandleFunc(http.MethodPost, "/dashboard/team/invitations/:token/accept", ds.AcceptInvitation)
	// app.HandleFunc(http.MethodPost, "/dashboard/team/invitations/:token/reject", ds.RejectInvitation)

	return app.mux
}

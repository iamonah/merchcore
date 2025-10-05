package router

import (
	"net/http"

	"github.com/IamOnah/storefronthq/internal/app/auth"

	"github.com/gorilla/mux"
)

func SetupRouter(us *auth.UserService) http.Handler {
	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", us.RegisterUser).Methods(http.MethodPost)
	auth.HandleFunc("/verify", us.ActivateUser).Methods(http.MethodPost)
	auth.HandleFunc("/signin", us.UserSignin).Methods(http.MethodPost)
	auth.HandleFunc("/signout", us.SignOut).Methods(http.MethodPost)
	auth.HandleFunc("/resend", us.ResendVerificationToken).Methods(http.MethodPost)
	auth.HandleFunc("/forgot-password", us.ForgotPassword).Methods(http.MethodPost)
	auth.HandleFunc("/reset-password", us.ResetPassword).Methods(http.MethodPost)
	auth.HandleFunc("/change-password", us.ChangePassword).Methods(http.MethodPost)
	auth.HandleFunc("/token/renew_access", us.RenewAccessToken).Methods(http.MethodPost)

	// tenant := r.PathPrefix("/store").Subrouter()
	// tenant.HandleFunc("/setup", SetupTenantHandler).Methods("POST")
	// tenant.HandleFunc("", GetTenantHandler).Methods("GET")
	// tenant.HandleFunc("/update", UpdateTenantHandler).Methods("PUT")
	// tenant.HandleFunc("/delete", DeleteTenantHandler).Methods("DELETE")
	// tenant.HandleFunc("/plan/upgrade", UpgradePlanHandler).Methods("POST")
	// tenant.HandleFunc("/plan/cancel", CancelPlanHandler).Methods("POST")

	// stores := r.PathPrefix("/stores").Subrouter()
	// stores.HandleFunc("", ListStoresHandler).Methods("GET")                 // list all stores
	// stores.HandleFunc("/create", CreateStoreHandler).Methods("POST")        // create new store
	// stores.HandleFunc("/{id}", GetStoreHandler).Methods("GET")              // get store details
	// stores.HandleFunc("/{id}/update", UpdateStoreHandler).Methods("PUT")    // update store info
	// stores.HandleFunc("/{id}/delete", DeleteStoreHandler).Methods("DELETE") // delete store

	// products := r.PathPrefix("/stores/{storeID}/products").Subrouter()
	// products.HandleFunc("", ListProductsHandler).Methods("GET")                 // list products
	// products.HandleFunc("/create", CreateProductHandler).Methods("POST")        // create product
	// products.HandleFunc("/{id}", GetProductHandler).Methods("GET")              // get product
	// products.HandleFunc("/{id}/update", UpdateProductHandler).Methods("PUT")    // update product
	// products.HandleFunc("/{id}/delete", DeleteProductHandler).Methods("DELETE") // delete product
	// products.HandleFunc("/{id}/stock", UpdateStockHandler).Methods("PATCH")     // update stock qty

	// orders := r.PathPrefix("/stores/{storeID}/orders").Subrouter()
	// orders.HandleFunc("", ListOrdersHandler).Methods("GET")                             // list orders
	// orders.HandleFunc("/{id}", GetOrderHandler).Methods("GET")                          // get order
	// orders.HandleFunc("/{id}/update-status", UpdateOrderStatusHandler).Methods("PATCH") // update status

	// customers := r.PathPrefix("/stores/{storeID}/customers").Subrouter()
	// customers.HandleFunc("", ListCustomersHandler).Methods("GET")                 // list customers
	// customers.HandleFunc("/{id}", GetCustomerHandler).Methods("GET")              // get customer
	// customers.HandleFunc("/{id}/delete", DeleteCustomerHandler).Methods("DELETE") // delete customer

	// billing := r.PathPrefix("/billing").Subrouter()
	// billing.HandleFunc("/subscribe", SubscribeHandler).Methods("POST")       // start subscription
	// billing.HandleFunc("/cancel", CancelSubscriptionHandler).Methods("POST") // cancel subscription
	// billing.HandleFunc("/invoices", ListInvoicesHandler).Methods("GET")      // list invoices

	// checkout := r.PathPrefix("/checkout").Subrouter()
	// checkout.HandleFunc("/{storeID}", CheckoutHandler).Methods("POST") // customer checkout

	// admin := r.PathPrefix("/admin").Subrouter()
	// admin.HandleFunc("/users", ListAllUsersHandler).Methods("GET")     // list all users
	// admin.HandleFunc("/tenants", ListAllTenantsHandler).Methods("GET") // list all tenants
	// admin.HandleFunc("/tenants/{id}", GetTenantHandler).Methods("GET") // get tenant detail
	// admin.HandleFunc("/stats", SystemStatsHandler).Methods("GET")      // system stats

	return r
}

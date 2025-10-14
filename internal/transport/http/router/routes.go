package router

import (
	"net/http"

	"github.com/IamOnah/storefronthq/internal/app/auth"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/middleware"
	"github.com/rs/zerolog"
)

func SetupRouter(us *auth.UserService, log *zerolog.Logger, maker *authz.JWTAuthMaker) http.Handler {
	app := NewApp(log, middleware.RecoverPanic(log))

	authbearer := middleware.AuthBearer(maker)
	// version := "v1"
	app.HandleFunc(http.MethodPost, "/auth/register", us.RegisterUser)
	app.HandleFunc(http.MethodPost, "/auth/activate", us.ActivateUser, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/signin", us.Authenticate)
	app.HandleFunc(http.MethodPost, "/auth/signout", us.SignOut, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/resend-token", us.ResendVerificationToken, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/forgot-password", us.ForgotPassword)
	app.HandleFunc(http.MethodPost, "/auth/reset-password", us.ResetPassword)
	app.HandleFunc(http.MethodPost, "/auth/change-password", us.ChangePassword, authbearer)
	app.HandleFunc(http.MethodPost, "/auth/token/renew_access", us.RenewAccessToken, authbearer)

	

	return app.mux
}

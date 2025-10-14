package router

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/IamOnah/storefronthq/internal/sdk/base"
	"github.com/IamOnah/storefronthq/internal/sdk/errs"
	"github.com/IamOnah/storefronthq/internal/sdk/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
)

type App struct {
	log      *zerolog.Logger
	mux      *mux.Router
	globalMw []base.Middleware
}

func NewApp(log *zerolog.Logger, globalMw ...base.Middleware) *App {
	return &App{
		log:      log,
		mux:      mux.NewRouter(),
		globalMw: globalMw,
	}
}

func (a *App) HandleFunc(method string, path string, handler base.HTTPHandlerWithErr, mw ...base.Middleware) {
	allMw := append(a.globalMw, mw...)
	wrapped := base.Chain(handler, allMw...)

	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIdKey, reqID))
		nwr := middleware.NewResponseWriter(w)
		start := time.Now()
		ip := base.GetClientIP(r)

		defer func() {
			if rec := recover(); rec != nil {
				a.log.Error().
					Interface("panic", rec).
					Bytes("stack", debug.Stack()).
					Msg("panic recovered in router")
				base.WriteJSONInternalError(nwr)
			}
			a.log.Info().
				Str("request_id", reqID).
				Str("method", r.Method).
				Str("url", r.URL.Path).
				Str("client_ip", ip).
				Str("user_agent", r.UserAgent()).
				Int("status_code", nwr.StatusCode).
				Dur("latency", time.Since(start)).
				Msg("incoming request")
		}()

		if err := wrapped(nwr, r); err != nil {
			a.handleError(nwr, r, err)
		}
	}).Methods(method)
}

func (app *App) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	reqID := r.Context().Value(middleware.RequestIdKey).(string)
	log := app.log.With().
		Str("req_id", reqID).
		Logger()

	if appErr, ok := err.(*errs.AppErr); ok {
		event := "http.request.failed"
		if appErr.Code >= 500 {
			log.Error().
				Err(fmt.Errorf("[internal]: %v", appErr)).
				Str("func_name", appErr.FuncName).
				Str("file_name", appErr.FileName).
				Str("event", event).
				Int("code", appErr.Code).Send()

			base.WriteJSONInternalError(w)
			return
		}

		base.WriteJSONError(w, appErr)
		return
	}

	log.Error().
		Err(err).
		Str("event", "http.request.unexpected").
		Msg("[panic] unhandled or non-AppErr error")

	base.WriteJSONInternalError(w)
}

package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/sdk/base"
	"github.com/iamonah/merchcore/internal/sdk/errs"
	"github.com/iamonah/merchcore/internal/sdk/midd"
	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
)

type App struct {
	log      *zerolog.Logger
	mux      *mux.Router
	globalMw []midd.Middleware
}

var RequestIDHeader = "X-Request-Id"

func NewApp(log *zerolog.Logger, globalMw ...midd.Middleware) *App {
	return &App{
		log:      log,
		mux:      mux.NewRouter(),
		globalMw: globalMw,
	}
}

func (a *App) HandleFunc(method string, path string, handler midd.HTTPHandlerWithErr, mw ...midd.Middleware) {
	allMw := append(a.globalMw, mw...)
	wrappedHandler := midd.Chain(handler, allMw...)

	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		w.Header().Set(RequestIDHeader, reqID)
		r = r.WithContext(context.WithValue(r.Context(), midd.RequestIdKey, reqID))
		nwr := midd.NewResponseWriter(w)
		start := time.Now()
		ip := base.GetClientIP(r)

		defer func() {
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

		if err := wrappedHandler(nwr, r); err != nil {
			a.handleError(nwr, r, err)
			return
		}
	}).Methods(method)
}

func (app *App) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	reqID := r.Context().Value(midd.RequestIdKey).(string)
	log := app.log.With().Str("req_id", reqID).Logger()

	//by design I always expect an *errs.AppErr
	if appErr, ok := err.(*errs.AppErr); ok {
		event := "http.request.failed"
		if appErr.Code >= http.StatusInternalServerError {
			log.Error().
				Err(fmt.Errorf("[internal]: %v", appErr)).
				Str("func_name", appErr.FuncName).
				Str("file_name", appErr.FileName).
				Str("event", event).
				Int("status_code", appErr.Code).Send()

			base.WriteJSONInternalError(w, appErr)
			return
		}

		base.WriteJSONError(w, appErr)
		return
	}

	base.WriteJSONInternalError(w, err)
}

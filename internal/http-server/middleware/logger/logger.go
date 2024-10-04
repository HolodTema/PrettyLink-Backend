package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// custom logger for all the requests
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		//create the copy of previous logger with caption of middleware
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		//print when this middleware is enabled, only 1 time
		log.Info("logger middleware enabled.")

		fn := func(w http.ResponseWriter, r *http.Request) {
			//this part executes before handling request:
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			//chi wrapper
			//it is needed to get information about response
			wrapResponseWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			//use timeNow to calculate how much time was spent to request handling
			timeNow := time.Now()
			//this defer will be executed after final request handling in ServeHttp(..)
			//so, we log the request after it was handled
			defer func() {
				entry.Info(
					"request completed.",
					slog.Int("status", wrapResponseWriter.Status()),
					slog.Int("bytes", wrapResponseWriter.BytesWritten()),
					slog.String("duration", time.Since(timeNow).String()),
				)
			}()
			//pass the control to the next middleware
			next.ServeHTTP(wrapResponseWriter, r)
		}
		//return this func as func, which can handle http requests
		return http.HandlerFunc(fn)
	}
}

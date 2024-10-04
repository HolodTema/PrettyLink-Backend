package main

import (
	"PrettyLinkBackend/internal/config"
	"PrettyLinkBackend/internal/http-server/handlers/url/save"
	"PrettyLinkBackend/internal/http-server/middleware/logger"
	"PrettyLinkBackend/internal/lib/logger/handlers/slogpretty"
	"PrettyLinkBackend/internal/lib/logger/sl"
	"PrettyLinkBackend/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting PrettyLink", slog.String("env", cfg.Env))
	log.Debug("Debug messages are working!")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()
	//connect middleware to the router
	//there are handlers sequence when some request is handled
	//the main handler, which is handles the request is the main handler
	//other handlers - are middlewares
	//for example, handler for authorization is middleware

	//this handler adds requestId to every request. It is good for debugging and tracing
	router.Use(middleware.RequestID) //middleware from chi lib
	router.Use(middleware.RealIP)    // to see user's ip, maybe it will be useful
	router.Use(middleware.Logger)    //this middleware logs all the requests automatically
	router.Use(logger.New(log))      //our custom middleware for logging requests
	router.Use(middleware.URLFormat) //to formalize api paths
	router.Use(middleware.Recoverer) // not to panic in case of failed request

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address)) // TODO: 01:38 stopped

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// srv.ListenAndServe() - is blocking call
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	//so it is not good if the row below will be executed. Because above was blocking call.
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	apiPkg "github.com/liuerfire/boxpractice/cmd/boxpractice/api"
	"github.com/liuerfire/boxpractice/pkg/httphandlers"
	"github.com/liuerfire/boxpractice/pkg/log"
	"github.com/liuerfire/boxpractice/pkg/store"
)

var (
	addr      = flag.String("addr", ":8080", "The addr to listen")
	verbosity = flag.Int("v", 0, "Number for the log level verbosity")
)

func main() {
	flag.Parse()

	logger := log.Init(*verbosity)

	setupLogger := logger.WithName("setup")

	sqlStore, err := store.NewSQLStore()
	if err != nil {
		setupLogger.Error(err, "failed to connect")
		os.Exit(1)
	}

	ctx := context.Background()

	router := mux.NewRouter()
	router.HandleFunc("/-/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	api, err := apiPkg.InitAPIHandler(ctx, logger, sqlStore)
	if err != nil {
		setupLogger.Error(err, "failed to connect")
		os.Exit(1)
	}
	api.RegisterRouter(router)

	handler := httphandlers.Register(
		router,
		httphandlers.LoggingHandler(logger.WithName("accesslog")),
		httphandlers.CorsConfigHandler(),
	)

	server := http.Server{
		Addr:    *addr,
		Handler: handler,
	}

	go func() {
		setupLogger.Info("start server")
		if err := server.ListenAndServe(); err != nil {
			setupLogger.Error(err, "failed to start server")
			os.Exit(1)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	server.Shutdown(ctx)
}

package app

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/orekhovskiy/shrtn/config"
	file "github.com/orekhovskiy/shrtn/internal/adapter/file/urlrepo"
	postgres "github.com/orekhovskiy/shrtn/internal/adapter/postgres/urlrepo"
	"github.com/orekhovskiy/shrtn/internal/handler/http"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api/shorten"
	"github.com/orekhovskiy/shrtn/internal/logger"
	"github.com/orekhovskiy/shrtn/internal/service/urlservice"
)

func Run(opts *config.Config) {
	zapLogger, err := logger.NewZapLogger() // Assuming logger is the package where ZapLogger is defined
	if err != nil {
		panic(fmt.Sprintf("unable to create zap logger: %v", err))
	}
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			panic(fmt.Sprintf("unable to sync zap logger: %v", err))
		}
	}()

	opts.LogConfig(zapLogger)

	var repo urlservice.Repository
	if opts.DatabaseDSN != "" {
		repo, err = postgres.NewRepository(*opts)
		if err != nil {
			panic(fmt.Sprintf("unable to connect to database: %v", err))
		}
	} else {
		fileRepo := file.NewRepository(*opts)
		err = fileRepo.LoadAll()
		repo = fileRepo
		if err != nil {
			panic(fmt.Sprintf("unable to load records from storage: %v", err))
		}
	}

	service := urlservice.NewService(repo)
	apiHandler := api.NewHandler(zapLogger, opts, service)
	shortenHandler := shorten.NewHandler(zapLogger, opts, service)

	router := http.NewRouter()
	router.
		WithHandler(apiHandler).
		WithHandler(shortenHandler)

	server := http.NewServer(opts)
	server.RegisterRoutes(router)

	zapLogger.Info("starting server",
		zap.String("address", opts.ServerAddress),
	)

	if err := server.Start(); err != nil {
		panic(fmt.Sprintf("unable to start a server: %v", err))
	}
}

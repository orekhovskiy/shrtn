package app

import (
	"fmt"

	"github.com/orekhovskiy/shrtn/internal/logger"

	"go.uber.org/zap"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/adapter/maprepo/urlrepo"
	"github.com/orekhovskiy/shrtn/internal/handler/http"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"github.com/orekhovskiy/shrtn/internal/service/urlservice"
)

func Run(opts *config.Config) {
	//logger, _ := zap.NewProduction()
	//defer func(logger *zap.Logger) {
	//	if err := logger.Sync(); err != nil {
	//		panic(fmt.Sprintf("unable to sync zap logger:%v", err))
	//	}
	//}(logger)

	zapLogger, err := logger.NewZapLogger() // Assuming logger is the package where ZapLogger is defined
	if err != nil {
		panic(fmt.Sprintf("unable to create zap logger: %v", err))
	}
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			panic(fmt.Sprintf("unable to sync zap logger: %v", err))
		}
	}()

	repo := urlrepo.NewRepository()
	service := urlservice.NewService(repo)
	handler := api.NewHandler(zapLogger, opts, *service)

	router := http.NewRouter()
	router.WithHandler(zapLogger, *handler)

	server := http.NewServer(opts)
	server.RegisterRoutes(router)

	zapLogger.Info("starting server",
		zap.String("address", opts.ServerAddress),
	)

	if err := server.Start(); err != nil {
		panic(fmt.Sprintf("unable to start a server: %v", err))
	}
}

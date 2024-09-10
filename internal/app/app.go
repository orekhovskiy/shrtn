package app

import (
	"fmt"
	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/adapter/maprepo/urlrepo"
	"github.com/orekhovskiy/shrtn/internal/handler/http"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"github.com/orekhovskiy/shrtn/internal/service/urlservice"
)

func Run(opts *config.Config) {
	repo := urlrepo.NewRepository()
	service := urlservice.NewService(repo)
	handler := api.NewHandler(opts, *service)

	router := http.NewRouter()
	router.WithHandler(*handler)

	server := http.NewServer(opts)
	server.RegisterRoutes(router)

	fmt.Printf("Starting server on %s\n", opts.ServerAddress)
	err := server.Start()
	if err != nil {
		fmt.Printf("unable to start a server: %s", err)
		return
	}
}

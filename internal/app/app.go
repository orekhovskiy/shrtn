package app

import (
	"log"

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

	log.Printf("starting server on %s", opts.ServerAddress)
	err := server.Start()
	if err != nil {
		log.Printf("unable to start a server: %s", err)
		return
	}
}

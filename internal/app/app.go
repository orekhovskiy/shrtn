package app

import (
	"fmt"
	"github.com/orekhovskiy/shrtn/internal/adapter/maprepo/urlrepo"
	"github.com/orekhovskiy/shrtn/internal/handler/http"
	"github.com/orekhovskiy/shrtn/internal/handler/http/api"
	"github.com/orekhovskiy/shrtn/internal/service/urlservice"
)

func Run() {
	repo := urlrepo.NewRepository()
	service := urlservice.NewService(repo)
	handler := api.NewHandler(*service)

	router := http.NewRouter()
	router.WithHandler(*handler)

	server := http.NewServer()
	server.RegisterRoutes(router)

	fmt.Println("Starting server on http://localhost:8080")
	err := server.Start()
	if err != nil {
		fmt.Printf("unable to start a server: %s", err)
		return
	}
}

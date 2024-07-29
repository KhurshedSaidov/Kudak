package main

import (
	"Kudak/configs"
	"Kudak/internal/handlers"
	"Kudak/internal/repository"
	"Kudak/internal/service"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	err := Run()
	if err != nil {
		return
	}
}

func Run() error {
	err := configs.InitDatabase()
	if err != nil {
		panic(err)
	}

	Repository := repository.NewRepository(configs.DB)
	Service := service.Service{Repository}
	Handler := handlers.Handler{&Service}

	router := InitRouters(&Handler)
	config, err := configs.InitConfigs()

	if err != nil {
		return err
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(router)

	address := config.Server.Host + ":" + config.Server.Port
	err = http.ListenAndServe(address, handler)

	if err != nil {
		log.Println("listen and serve error", err)
		return err
	}
	return nil
}

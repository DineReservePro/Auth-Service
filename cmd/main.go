package main

import (
	"auth-service/api"
	"auth-service/api/handler"
	"auth-service/config"
	"auth-service/server"
	"auth-service/storage/postgres"
	"log"
	"sync"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg := config.Load()
	router := api.Routes(handler.NewHandler(db))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		server.ServerRun(db)
	}()

	go func() {
		defer wg.Done()
		log.Fatal(router.Run(cfg.HTTP_PORT))
	}()

	wg.Wait()
}

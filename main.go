package main

import (
	"log"
	"net/http"
	"time"
	"user-service/config"
	"user-service/db"
	"user-service/server"
	"user-service/services"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	authService := services.NewAuthService(authRepo, conf)

	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
		AuthService:    authService,
	}
	s.Start()
}


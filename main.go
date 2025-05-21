package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aakash-Pandit/reetro-golang/common"
	"github.com/Aakash-Pandit/reetro-golang/config"
	"github.com/Aakash-Pandit/reetro-golang/routes"
	"github.com/Aakash-Pandit/reetro-golang/server"
	"github.com/Aakash-Pandit/reetro-golang/storages"
	"github.com/gorilla/mux"
)

func init() {
	err := config.ReadEnvironmentVariables()
	if err != nil {
		log.Fatal("Unable to read environment variables:", err)
		return
	}
}

func main() {
	database, err := storages.DBInit()
	if err != nil {
		log.Fatal("Unable to connect DB:", err)
		return
	}

	email := &common.Email{
		EmailFrom: os.Getenv("EMAIL_ID"),
		Password:  os.Getenv("EMAIL_PASSWORD"),
		Host:      os.Getenv("EMAIL_HOST"),
		Address:   os.Getenv("EMAIL_HOST") + ":" + os.Getenv("EMAIL_PORT"),
	}

	defer database.Close()

	log.Println("Application has been started on :", os.Getenv("APPLICATION_PORT"))

	route := routes.NewRouter(
		mux.NewRouter(),
		fmt.Sprintf(":%s", os.Getenv("APPLICATION_PORT")),
		database.Postgres,
		database.Redis,
		email,
	)

	server := server.NewServer(route)
	server.Start()
}

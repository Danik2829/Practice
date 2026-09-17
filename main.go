package main

import (
	"fmt"
	"log"
	"net/http"
	"restAPI/config"
	"restAPI/handlers"
	"restAPI/service"
	"restAPI/storage"

	"github.com/gorilla/mux"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		return
	}

	store, err := storage.NewDbStorage(*conf)
	log.Printf("MAIN: store=%p", store)

	if err != nil {
		log.Printf("MAIN: NewDbStorage error: %v", err)
		return
	}
	service := service.NewUserService(store)
	handler := handlers.NewHandler(service)

	mux := mux.NewRouter()
	handler.InitRoutes(mux)

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(port, mux)
}

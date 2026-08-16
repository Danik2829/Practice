package main

import (
	"fmt"
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
	if err != nil {
		//defer store.Close()
		//return
	}
	service := service.NewUserService(store)
	handler := handlers.NewHandler(service)

	mux := mux.NewRouter()
	handler.InitRoutes(mux)

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(port, mux)
}

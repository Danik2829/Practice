package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"restAPI/handlers"
	"restAPI/service"
	"restAPI/storage"
)

func main() {
	stor, err := storage.NewDbStorage()
	if err != nil {
		fmt.Println(err)
		defer stor.Close()
		return
	}
	service := service.NewUserService(stor)
	handler := handlers.NewHandler(service)

	mux := mux.NewRouter()
	mux.HandleFunc("/users", handler.CreateUser).Methods("POST")
	mux.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	mux.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	mux.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(port, mux)
}

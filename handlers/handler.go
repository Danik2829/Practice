package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"restAPI/core"

	"github.com/gorilla/mux"
)

type Handler struct {
	service core.UserService
}

func NewHandler(s core.UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) InitRoutes(mux *mux.Router) {
	mux.HandleFunc("/users", h.CreateUser).Methods("POST")
	mux.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	mux.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	mux.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.service.Create(user)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if encodeErr := json.NewEncoder(w).Encode(user); encodeErr != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
		return
	}
	switch {
	case errors.Is(err, core.InvalidData):
		http.Error(w, "Invalid user data", http.StatusUnprocessableEntity)
	case errors.Is(err, core.UserExist):
		http.Error(w, "User already exists", http.StatusConflict)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	if user := h.service.Get(id["id"]); user == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else {
		rawJson, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rawJson)
		return
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user.ID = id
	updatedUser, err := h.service.Update(user)
	switch {
	case errors.Is(err, core.NotFound):
		http.Error(w, "User not found", http.StatusNotFound)
		return
	case errors.Is(err, core.InvalidData):
		http.Error(w, "Invalid user data", http.StatusUnprocessableEntity)
		return
	case err != nil:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idURL := mux.Vars(r)
	if id := h.service.Delete(idURL["id"]); id == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

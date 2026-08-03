package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"restAPI/core"

	"github.com/gorilla/mux"
)

type Handler struct {
	service core.UserServiceInterface
}

func NewHandler(s core.UserServiceInterface) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.service.CreateUser(user)
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
	if user, err := h.service.GetUser(id["id"]); err == core.NotFound {
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
	idURL := mux.Vars(r)
	var user core.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err == nil {
		user.ID = idURL["id"]
		jsonErr := h.service.UpdateUser(user)
		switch jsonErr {
		case core.NotFound:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		case core.InvalidData:
			http.Error(w, "JSON parsed but values are wrong", http.StatusUnprocessableEntity)
			return
		default:
			rawJson, _ := json.Marshal(user)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(rawJson)
			return
		}
	} else {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idURL := mux.Vars(r)
	if err := h.service.DeleteUser(idURL["id"]); err == core.NotFound {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

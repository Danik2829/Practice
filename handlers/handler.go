package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restAPI/core"
	"restAPI/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.UserService
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user core.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err == nil {
		jsonErr := h.service.CreateUser(user)
		fmt.Println(jsonErr)
		switch jsonErr {
		case core.InvalidData:
			http.Error(w, "JSON parsed but values are wrong", http.StatusUnprocessableEntity)
			return
		case core.UserExist:
			http.Error(w, "User with this ID already exist", http.StatusConflict)
			return
		default:
			rawJson, _ := json.Marshal(user)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(rawJson)
			return
		}
	} else {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
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

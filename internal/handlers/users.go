package handlers

import (
	"encoding/json"
	"fmt"
	"golang-firstcode/internal/models"
	"golang-firstcode/internal/service"
	"golang-firstcode/pkg/common"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	common.WriteJSON(w, http.StatusOK, users, "Successfully fetched users")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", id) {
			common.WriteError(w, http.StatusNotFound, "User not found")
		} else {
			common.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	common.WriteJSON(w, http.StatusOK, user, "Successfully fetched user")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Name == "" {
		common.WriteError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	createdUser, err := h.service.Create(user)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	common.WriteJSON(w, http.StatusCreated, createdUser, "Successfully created user")
}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Name == "" {
		common.WriteError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	user.ID = id
	updatedUser, err := h.service.Update(user)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with ID %d not found", id) {
			common.WriteError(w, http.StatusNotFound, "User not found")
		} else {
			common.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	common.WriteJSON(w, http.StatusOK, updatedUser, "Successfully updated user")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		common.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	common.WriteJSON(w, http.StatusOK, nil, "User deleted successfully")
}

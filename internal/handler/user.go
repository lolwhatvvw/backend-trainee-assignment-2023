package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/models"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/storage"
)

const (
	ErrInvalidUserID    = "invalid user ID"
	ErrMismatchedUserID = "mismatched user ID"
	ErrMissingUserID    = "missing user ID"
)

type UserHandler struct {
	us storage.UserStorage
}

func NewUserHandler(us storage.UserStorage) *UserHandler {
	return &UserHandler{us: us}
}

// ListUsers godoc
// @Summary List all users
// @Description Returns a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.us.GetUsers()
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	if len(users) == 0 {
		_ = render.Render(w, r, ErrNotFound())
		return
	}

	render.JSON(w, r, users)
}

type CreateUserRequest struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Username  string `json:"username" validate:"required"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "The user to create"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.us.CreateUser(&user); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// ReadUser godoc
// @Summary Get a user
// @Description Returns a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID of the user to retrieve"
// @Success 200 {object} models.User
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrInvalidUserID)))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrInvalidUserID)))
		return
	}

	user, err := h.us.GetUserByID(id)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if user == nil {
		render.Render(w, r, ErrNotFound())
		return
	}

	render.JSON(w, r, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Updates an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID of the user to update"
// @Param user body models.User true "The user data to update"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrMissingUserID)))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrInvalidUserID)))
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if user.ID != id {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrMismatchedUserID)))
		return
	}

	if err := h.us.UpdateUser(&user); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID of the user to delete"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrMissingUserID)))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf(ErrInvalidUserID)))
		return
	}

	if err := h.us.DeleteUser(id); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

type updateUserSegments struct {
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToRemove []string `json:"segments_to_remove"`
}

// UpdateUserSegments godoc
//
// @Summary Update the segments of a user
// @Description Updates the segments of an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID of the user to update segments for"
// @Param update body updateUserSegments true "The segments to add or remove"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id}/segments [put]
func (h *UserHandler) UpdateUserSegments(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		render.Render(w, r, ErrRender(fmt.Errorf(ErrMissingUserID)))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("invalid user ID parameter: %v\n", idStr)
		render.Render(w, r, ErrRender(fmt.Errorf(ErrInvalidUserID)))
		return
	}

	var update updateUserSegments
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("failed to decode user data from request: %v\n", err)
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid request body")))
		return
	}

	if err := h.us.UpdateUserSegments(id, update.SegmentsToAdd, update.SegmentsToRemove); err != nil {
		log.Printf("failed to update user segments: %v\n", err)
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

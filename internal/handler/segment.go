package handler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/models"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/storage"
)

type SegmentHandler struct {
	ss storage.SegmentStorage
}

func NewSegmentHandler(ss storage.SegmentStorage) *SegmentHandler {
	return &SegmentHandler{ss: ss}
}

// ListSegments godoc
//
// @Summary List all segments
// @Description Returns a list of all segments in the system
// @Tags segments
// @Accept json
// @Produce json
// @Success 200 {array} models.Segment
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/segments [get]
func (h *SegmentHandler) ListSegments(w http.ResponseWriter, r *http.Request) {
	segments, err := h.ss.GetSegments()
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.JSON(w, r, segments)
}

type CreateSegmentRequest struct {
	Name string `json:"name"`
}

// CreateSegment godoc
//
// @Summary Create a new segment
// @Description Creates a new segment in the system
// @Tags segments
// @Accept json
// @Produce json
// @Param segment body CreateSegmentRequest true "The segment to create"
// @Success 201 {object} models.Segment
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/segments [post]
func (h *SegmentHandler) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var segment models.Segment
	if err := json.NewDecoder(r.Body).Decode(&segment); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.ss.CreateSegment(&segment); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, segment)
}

// ReadSegment godoc
//
// @Summary Get a segment
// @Description Returns a single segment by slug
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the segment to retrieve"
// @Success 200 {object} models.Segment
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/segments/{slug} [get]
func (h *SegmentHandler) ReadSegment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Render(w, r, ErrMissingField("slug"))
		return
	}

	segment, err := h.ss.GetSegmentByName(slug)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if segment == nil {
		render.Render(w, r, ErrNotFound())
		return
	}

	render.JSON(w, r, segment)
}

// UpdateSegment godoc
//
// @Summary Update a segment
// @Description Updates an existing segment by slug
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the segment to update"
// @Param segment body models.Segment true "The segment data to update"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/segments/{slug} [put]
func (h *SegmentHandler) UpdateSegment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Render(w, r, ErrMissingField("slug"))
		return
	}

	var segment models.Segment
	if err := json.NewDecoder(r.Body).Decode(&segment); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.ss.UpdateSegment(&segment); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// DeleteSegment godoc
//
// @Summary Delete a segment
// @Description Deletes an existing segment by slug
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the segment to delete"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/segments/{slug} [delete]
func (h *SegmentHandler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Render(w, r, ErrMissingField("slug"))
		return
	}

	if err := h.ss.DeleteSegmentBySlug(slug); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// ListUsersInSegment godoc
//
// @Summary List all users in a segment
// @Description Returns a list of all users in the specified segment
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the segment to retrieve users for"
// @Success 200 {array} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/segments/{slug}/users [get]
func (h *SegmentHandler) ListUsersInSegment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Render(w, r, ErrMissingField("slug"))
		return
	}

	users, err := h.ss.GetUsersInSegment(slug)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.JSON(w, r, users)
}

// AddUserToSegment godoc
//
// @Summary Add a user to a segment
// @Description Adds a user to the specified segment
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the segment to add the user to"
// @Param id path int true "ID of the user to add to the segment"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/segments/{slug}/users/{id} [put]
func (h *SegmentHandler) AddUserToSegment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Render(w, r, ErrMissingField("slug"))
		return
	}

	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidField("id", userIDStr))
		return
	}

	if err := h.ss.AddUserToSegment(slug, userID); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

// DeleteUserFromSegment godoc
//
// @Summary Remove a user from a segment
// @Description Removes a user from the specified segment
// @Tags segments
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the segment to remove the user from"
// @Param id path int true "ID of the user to remove from the segment"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/segments/{slug}/users/{id} [delete]
func (h *SegmentHandler) DeleteUserFromSegment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		render.Render(w, r, ErrMissingField("slug"))
		return
	}

	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidField("id", userIDStr))
		return
	}

	if err := h.ss.DeleteUserFromSegment(slug, userID); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

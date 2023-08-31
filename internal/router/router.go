package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/handler"
	"net/http"
	"time"
)

const TIMEOUT = 60 * time.Second

func GetRouter(userController *handler.UserHandler, segmentController *handler.SegmentHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(TIMEOUT))

	buildTree(r, userController, segmentController)

	return r
}

func buildTree(r *chi.Mux, userController *handler.UserHandler, segmentController *handler.SegmentHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", userRouter(userController))
		r.Mount("/segments", segmentRouter(segmentController))
	})
}

func userRouter(userController *handler.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", userController.ListUsers)
	r.Post("/", userController.CreateUser)
	r.Get("/{id}", userController.ReadUser)
	r.Put("/{id}", userController.UpdateUser)
	r.Delete("/{id}", userController.DeleteUser)
	r.Put("/{id}/segments", userController.UpdateUserSegments)
	return r
}

func segmentRouter(segmentController *handler.SegmentHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", segmentController.ListSegments)
	r.Post("/", segmentController.CreateSegment)
	r.Get("/{slug}", segmentController.ReadSegment)
	r.Put("/{slug}", segmentController.UpdateSegment)
	r.Delete("/{slug}", segmentController.DeleteSegment)
	r.Get("/{slug}/users", segmentController.ListUsersInSegment)
	r.Put("/{slug}/users/{id}", segmentController.AddUserToSegment)
	r.Delete("/{slug}/users/{id}", segmentController.DeleteUserFromSegment)
	return r
}

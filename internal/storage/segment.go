package storage

import "github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/models"

type SegmentStorage interface {
	CreateSegment(segment *models.Segment) error
	GetSegments() ([]*models.Segment, error)
	GetSegmentByName(slug string) (*models.Segment, error)
	UpdateSegment(segment *models.Segment) error
	DeleteSegmentBySlug(slug string) error
	GetUsersInSegment(slug string) ([]*models.User, error)
	AddUserToSegment(slug string, userID int64) error
	DeleteUserFromSegment(slug string, userID int64) error
}

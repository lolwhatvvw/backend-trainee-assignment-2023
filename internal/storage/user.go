package storage

import "github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/models"

type UserStorage interface {
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	GetUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
	UpdateUserSegments(id int64, segmentsToAdd, segmentsToRemove []string) error
}

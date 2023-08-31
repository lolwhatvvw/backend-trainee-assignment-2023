package postgres

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/models"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/storage"
	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) (storage.UserStorage, error) {
	return &userStorage{db: db}, nil
}

func (s *userStorage) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *userStorage) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	result := s.db.Preload("Segments").First(user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // no user found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}
	return user, nil
}

func (s *userStorage) GetUsers() ([]*models.User, error) {
	var users []*models.User
	result := s.db.Preload("Segments").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %w", result.Error)
	}
	return users, nil
}

func (s *userStorage) UpdateUser(user *models.User) error {
	result := s.db.Model(user).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected when updating user")
	}
	return nil
}

func (s *userStorage) DeleteUser(id int64) error {
	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected when deleting user")
	}
	return nil
}

func (s *userStorage) UpdateUserSegments(id int64, segmentsToAdd, segmentsToRemove []string) error {
	tx := s.db.Begin()

	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer tx.Rollback()

	// Get existing segments for user
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	segmentsToAddSet := make(map[string]bool)
	for _, segment := range segmentsToAdd {
		segmentsToAddSet[segment] = true
	}

	for _, segment := range user.Segments {
		delete(segmentsToAddSet, segment.Name)
	}

	segmentsToRemoveSet := make(map[string]bool)
	for _, segment := range segmentsToRemove {
		segmentsToRemoveSet[segment] = true
	}

	// Find segments in both sets (i.e., intersection)
	var intersection []string
	for segment := range segmentsToAddSet {
		if segmentsToRemoveSet[segment] {
			intersection = append(intersection, segment)
		}
	}

	// Remove segments in intersection from both sets
	for _, segment := range intersection {
		delete(segmentsToAddSet, segment)
		delete(segmentsToRemoveSet, segment)
	}

	// Add user to non-intersecting segments using single query
	if len(segmentsToAddSet) > 0 {
		err := s.bulkInsertUnique(id, segmentsToAddSet, tx)
		if err != nil {
			return err
		}
	}

	// Remove user from non-intersecting segments using single query
	if len(segmentsToRemoveSet) > 0 {
		err := s.bulkDeleteUnique(id, segmentsToRemoveSet, tx)
		if err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (s *userStorage) bulkDeleteUnique(id int64, segmentsToRemoveSet map[string]bool, tx *gorm.DB) error {
	segmentNames := make([]string, 0, len(segmentsToRemoveSet))
	values := make([]interface{}, 0, 1+len(segmentsToRemoveSet))
	values = append(values, id)
	idx := 2
	for segment := range segmentsToRemoveSet {
		segmentNames = append(segmentNames, fmt.Sprintf("$%d", idx))
		values = append(values, segment)
		idx++
	}
	query := fmt.Sprintf(
		"DELETE FROM user_segments WHERE user_id = $1 AND segment_name IN (%s)",
		strings.Join(segmentNames, ","),
	)
	result := tx.Exec(query, values...)
	if result.Error != nil {
		return fmt.Errorf("failed to remove user from segments: %w", result.Error)
	}
	return nil
}

func (s *userStorage) bulkInsertUnique(id int64, segmentsToAddSet map[string]bool, tx *gorm.DB) error {
	valueStrings := make([]string, 0, len(segmentsToAddSet))
	valueArgs := make([]any, 0, len(segmentsToAddSet)*2)
	i := 0
	for segment := range segmentsToAddSet {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, id)
		valueArgs = append(valueArgs, segment)
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO user_segments (user_id, segment_name) VALUES %s", strings.Join(valueStrings, ","),
	)
	log.Println(query)
	result := tx.Exec(query, valueArgs...)
	if result.Error != nil {
		return fmt.Errorf("failed to add user to segments: %w", result.Error)
	}
	return nil
}

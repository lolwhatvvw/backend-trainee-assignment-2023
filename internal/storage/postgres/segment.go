package postgres

import (
	"errors"
	"fmt"

	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/models"
	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/storage"
	"gorm.io/gorm"
)

type segmentStorage struct {
	db          *gorm.DB
	userStorage storage.UserStorage
}

func NewSegmentStorage(db *gorm.DB, us storage.UserStorage) (storage.SegmentStorage, error) {
	return &segmentStorage{db: db, userStorage: us}, nil
}

func (s *segmentStorage) CreateSegment(segment *models.Segment) error {
	return s.db.Create(segment).Error
}

func (s *segmentStorage) GetSegmentByName(name string) (*models.Segment, error) {
	segment := &models.Segment{}
	if err := s.db.Preload("Users").Where("name = ?", name).First(segment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("segment with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get segment by name '%s': %w", name, err)
	}
	return segment, nil
}

func (s *segmentStorage) GetSegments() ([]*models.Segment, error) {
	var segments []*models.Segment
	result := s.db.Preload("Users").Find(&segments)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get segments: %w", result.Error)
	}

	return segments, nil
}

func (s *segmentStorage) UpdateSegment(segment *models.Segment) error {
	result := s.db.Updates(segment)
	if result.Error != nil {
		return fmt.Errorf("failed to update segment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected when updating segment")
	}

	return nil
}

func (s *segmentStorage) DeleteSegmentBySlug(slug string) error {
	result := s.db.Where("slug = ?", slug).Delete(&models.Segment{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete segment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected when deleting segment")
	}

	return nil
}

func (s *segmentStorage) GetUsersInSegment(slug string) ([]*models.User, error) {
	var users []*models.User
	err := s.db.Model(&models.Segment{}).Where("name = ?", slug).Association("Users").Find(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get users in segment: %w", err)
	}

	return users, nil
}

func (s *segmentStorage) AddUserToSegment(slug string, userID int64) error {
	segment := &models.Segment{}
	if err := s.db.Where("name = ?", slug).First(segment).Error; err != nil {
		return fmt.Errorf("failed to get segment by name: %w", err)
	}

	user := &models.User{}
	if err := s.db.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", userID)
		}
		return fmt.Errorf("failed to get user by ID %d: %w", userID, err)
	}

	if err := s.db.Model(segment).Association("Users").Append(user); err != nil {
		return fmt.Errorf("failed to add user to segment: %w", err)
	}

	return nil
}

func (s *segmentStorage) DeleteUserFromSegment(slug string, userID int64) error {
	segment := &models.Segment{}
	if err := s.db.Where("name = ?", slug).Preload("Users").First(segment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("segment with name '%s' not found", slug)
		}
		return fmt.Errorf("failed to get segment by slug: %w", err)
	}

	user := &models.User{ID: userID}
	if err := s.db.First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", userID)
		}
		return fmt.Errorf("failed to get user by ID %d: %w", userID, err)
	}

	return s.db.Model(segment).Association("Users").Delete(user)
}

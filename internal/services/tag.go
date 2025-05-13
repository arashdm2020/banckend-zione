package services

import (
	"errors"
	"strings"

	"zionechainapi/internal/database"
	"zionechainapi/internal/models"
	"gorm.io/gorm"
)

// TagService handles tag-related operations
type TagService struct{}

// NewTagService creates a new tag service
func NewTagService() *TagService {
	return &TagService{}
}

// TagRequest represents the tag request
type TagRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(req TagRequest) (*TagResponse, error) {
	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if slug already exists
	var count int64
	if err := database.DB.Model(&models.Tag{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("tag with this name already exists")
	}

	// Create tag
	tag := models.Tag{
		Name: req.Name,
		Slug: slug,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

// UpdateTag updates a tag
func (s *TagService) UpdateTag(id uint, req TagRequest) (*TagResponse, error) {
	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if slug already exists and is not this tag
	var count int64
	if err := database.DB.Model(&models.Tag{}).Where("slug = ? AND id != ?", slug, id).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("tag with this name already exists")
	}

	// Update tag
	tag.Name = req.Name
	tag.Slug = slug

	if err := database.DB.Save(&tag).Error; err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(id uint) error {
	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tag not found")
		}
		return err
	}

	// Start transaction
	tx := database.DB.Begin()

	// Remove project associations
	if err := tx.Model(&tag).Association("Projects").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Remove blog associations
	if err := tx.Model(&tag).Association("BlogPosts").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Delete tag
	if err := tx.Delete(&tag).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// ListTags lists all tags
func (s *TagService) ListTags() ([]TagResponse, error) {
	var tags []models.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		return nil, err
	}

	var response []TagResponse
	for _, tag := range tags {
		response = append(response, TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response, nil
}

// GetTagByID gets a tag by ID
func (s *TagService) GetTagByID(id uint) (*TagResponse, error) {
	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return &TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
}

// GetTagBySlug gets a tag by slug
func (s *TagService) GetTagBySlug(slug string) (*TagResponse, error) {
	var tag models.Tag
	if err := database.DB.Where("slug = ?", slug).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return &TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}, nil
} 
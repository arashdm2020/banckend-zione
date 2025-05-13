package services

import (
	"errors"
	"strings"

	"zionechainapi/internal/database"
	"zionechainapi/internal/models"
	"gorm.io/gorm"
)

// CategoryService handles category-related operations
type CategoryService struct{}

// NewCategoryService creates a new category service
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// CategoryRequest represents the category request
type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// CategoryType represents the type of category
type CategoryType string

const (
	// CategoryTypeProject represents a project category
	CategoryTypeProject CategoryType = "project"
	// CategoryTypeBlog represents a blog category
	CategoryTypeBlog CategoryType = "blog"
)

// CreateProjectCategory creates a new project category
func (s *CategoryService) CreateProjectCategory(req CategoryRequest) (*ProjectCategoryResponse, error) {
	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if slug already exists
	var count int64
	if err := database.DB.Model(&models.ProjectCategory{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("category with this name already exists")
	}

	// Create category
	category := models.ProjectCategory{
		Name: req.Name,
		Slug: slug,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return &ProjectCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// UpdateProjectCategory updates a project category
func (s *CategoryService) UpdateProjectCategory(id uint, req CategoryRequest) (*ProjectCategoryResponse, error) {
	var category models.ProjectCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if slug already exists and is not this category
	var count int64
	if err := database.DB.Model(&models.ProjectCategory{}).Where("slug = ? AND id != ?", slug, id).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("category with this name already exists")
	}

	// Update category
	category.Name = req.Name
	category.Slug = slug

	if err := database.DB.Save(&category).Error; err != nil {
		return nil, err
	}

	return &ProjectCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// DeleteProjectCategory deletes a project category
func (s *CategoryService) DeleteProjectCategory(id uint) error {
	var category models.ProjectCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if category is used by any project
	var count int64
	if err := database.DB.Model(&models.Project{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("category is used by projects and cannot be deleted")
	}

	return database.DB.Delete(&category).Error
}

// ListProjectCategories lists all project categories
func (s *CategoryService) ListProjectCategories() ([]ProjectCategoryResponse, error) {
	var categories []models.ProjectCategory
	if err := database.DB.Find(&categories).Error; err != nil {
		return nil, err
	}

	var response []ProjectCategoryResponse
	for _, category := range categories {
		response = append(response, ProjectCategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return response, nil
}

// GetProjectCategoryByID gets a project category by ID
func (s *CategoryService) GetProjectCategoryByID(id uint) (*ProjectCategoryResponse, error) {
	var category models.ProjectCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &ProjectCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// CreateBlogCategory creates a new blog category
func (s *CategoryService) CreateBlogCategory(req CategoryRequest) (*BlogCategoryResponse, error) {
	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if slug already exists
	var count int64
	if err := database.DB.Model(&models.BlogCategory{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("category with this name already exists")
	}

	// Create category
	category := models.BlogCategory{
		Name: req.Name,
		Slug: slug,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return &BlogCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// UpdateBlogCategory updates a blog category
func (s *CategoryService) UpdateBlogCategory(id uint, req CategoryRequest) (*BlogCategoryResponse, error) {
	var category models.BlogCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if slug already exists and is not this category
	var count int64
	if err := database.DB.Model(&models.BlogCategory{}).Where("slug = ? AND id != ?", slug, id).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("category with this name already exists")
	}

	// Update category
	category.Name = req.Name
	category.Slug = slug

	if err := database.DB.Save(&category).Error; err != nil {
		return nil, err
	}

	return &BlogCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
}

// DeleteBlogCategory deletes a blog category
func (s *CategoryService) DeleteBlogCategory(id uint) error {
	var category models.BlogCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if category is used by any blog post
	var count int64
	if err := database.DB.Model(&models.BlogPost{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("category is used by blog posts and cannot be deleted")
	}

	return database.DB.Delete(&category).Error
}

// ListBlogCategories lists all blog categories
func (s *CategoryService) ListBlogCategories() ([]BlogCategoryResponse, error) {
	var categories []models.BlogCategory
	if err := database.DB.Find(&categories).Error; err != nil {
		return nil, err
	}

	var response []BlogCategoryResponse
	for _, category := range categories {
		response = append(response, BlogCategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return response, nil
}

// GetBlogCategoryByID gets a blog category by ID
func (s *CategoryService) GetBlogCategoryByID(id uint) (*BlogCategoryResponse, error) {
	var category models.BlogCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &BlogCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}, nil
} 
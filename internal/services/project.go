package services

import (
	"errors"
	"strings"
	"time"

	"zionechainapi/internal/database"
	"zionechainapi/internal/models"
	"gorm.io/gorm"
)

// ProjectService handles project-related operations
type ProjectService struct{}

// NewProjectService creates a new project service
func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// CreateProjectRequest represents the create project request
type CreateProjectRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	CategoryID  uint     `json:"category_id" binding:"required"`
	TagIDs      []uint   `json:"tag_ids"`
	Featured    bool     `json:"featured"`
	Published   bool     `json:"published"`
}

// UpdateProjectRequest represents the update project request
type UpdateProjectRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	CategoryID  uint    `json:"category_id"`
	TagIDs      []uint  `json:"tag_ids"`
	Featured    *bool   `json:"featured"`
	Published   *bool   `json:"published"`
}

// ProjectMediaRequest represents the project media request
type ProjectMediaRequest struct {
	Type      string `json:"type" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

// ProjectResponse represents the project response
type ProjectResponse struct {
	ID          uint                   `json:"id"`
	Title       string                 `json:"title"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	CategoryID  uint                   `json:"category_id"`
	Category    ProjectCategoryResponse `json:"category"`
	Media       []ProjectMediaResponse  `json:"media"`
	Tags        []TagResponse           `json:"tags"`
	Featured    bool                   `json:"featured"`
	Published   bool                   `json:"published"`
	CreatedBy   uint                   `json:"created_by"`
	UpdatedBy   uint                   `json:"updated_by"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// ProjectCategoryResponse represents the project category response
type ProjectCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProjectMediaResponse represents the project media response
type ProjectMediaResponse struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

// TagResponse represents the tag response
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req CreateProjectRequest, userID uint) (*ProjectResponse, error) {
	// Create slug from title
	slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))

	// Check if slug already exists
	var count int64
	if err := database.DB.Model(&models.Project{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		// Append timestamp to slug to make it unique
		slug = slug + "-" + string(time.Now().Unix())
	}

	// Create project
	project := models.Project{
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Content:     req.Content,
		CategoryID:  req.CategoryID,
		Featured:    req.Featured,
		Published:   req.Published,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	// Start transaction
	tx := database.DB.Begin()
	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add tags if any
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Model(&project).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load project with relationships
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").First(&project, project.ID).Error; err != nil {
		return nil, err
	}

	// Map to response
	return s.mapProjectToResponse(project), nil
}

// GetProjectByID gets a project by ID
func (s *ProjectService) GetProjectByID(id uint) (*ProjectResponse, error) {
	var project models.Project
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return s.mapProjectToResponse(project), nil
}

// GetProjectBySlug gets a project by slug
func (s *ProjectService) GetProjectBySlug(slug string) (*ProjectResponse, error) {
	var project models.Project
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").Where("slug = ?", slug).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return s.mapProjectToResponse(project), nil
}

// ListProjects lists all projects with pagination
func (s *ProjectService) ListProjects(page, limit int, categoryID uint, featured, published bool) ([]ProjectResponse, int64, error) {
	var projects []models.Project
	var total int64

	// Base query
	query := database.DB.Model(&models.Project{})

	// Apply filters
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if featured {
		query = query.Where("featured = ?", featured)
	}

	// Default to published only
	query = query.Where("published = ?", published)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * limit
	if err := query.Preload("Category").Preload("Media").Preload("Tags").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	// Map to response
	var response []ProjectResponse
	for _, project := range projects {
		response = append(response, *s.mapProjectToResponse(project))
	}

	return response, total, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(id uint, req UpdateProjectRequest, userID uint) (*ProjectResponse, error) {
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	// Update fields if provided
	tx := database.DB.Begin()

	if req.Title != "" && req.Title != project.Title {
		// Create new slug from title
		slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))

		// Check if slug already exists and is not this project
		var count int64
		if err := tx.Model(&models.Project{}).Where("slug = ? AND id != ?", slug, id).Count(&count).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if count > 0 {
			// Append timestamp to slug to make it unique
			slug = slug + "-" + string(time.Now().Unix())
		}

		project.Title = req.Title
		project.Slug = slug
	}

	if req.Description != "" {
		project.Description = req.Description
	}

	if req.Content != "" {
		project.Content = req.Content
	}

	if req.CategoryID > 0 {
		project.CategoryID = req.CategoryID
	}

	if req.Featured != nil {
		project.Featured = *req.Featured
	}

	if req.Published != nil {
		project.Published = *req.Published
	}

	project.UpdatedBy = userID

	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update tags if provided
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Model(&project).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load project with relationships
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").First(&project, id).Error; err != nil {
		return nil, err
	}

	return s.mapProjectToResponse(project), nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id uint) error {
	var project models.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return err
	}

	// Start transaction
	tx := database.DB.Begin()

	// Delete media
	if err := tx.Where("project_id = ?", id).Delete(&models.ProjectMedia{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Remove tag associations
	if err := tx.Model(&project).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Delete project
	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// AddProjectMedia adds media to a project
func (s *ProjectService) AddProjectMedia(projectID uint, req ProjectMediaRequest) (*ProjectMediaResponse, error) {
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	media := models.ProjectMedia{
		ProjectID: projectID,
		Type:      req.Type,
		URL:       req.URL,
		Caption:   req.Caption,
		SortOrder: req.SortOrder,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		return nil, err
	}

	return &ProjectMediaResponse{
		ID:        media.ID,
		Type:      media.Type,
		URL:       media.URL,
		Caption:   media.Caption,
		SortOrder: media.SortOrder,
	}, nil
}

// UpdateProjectMedia updates project media
func (s *ProjectService) UpdateProjectMedia(mediaID uint, req ProjectMediaRequest) (*ProjectMediaResponse, error) {
	var media models.ProjectMedia
	if err := database.DB.First(&media, mediaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("media not found")
		}
		return nil, err
	}

	media.Type = req.Type
	media.URL = req.URL
	media.Caption = req.Caption
	media.SortOrder = req.SortOrder

	if err := database.DB.Save(&media).Error; err != nil {
		return nil, err
	}

	return &ProjectMediaResponse{
		ID:        media.ID,
		Type:      media.Type,
		URL:       media.URL,
		Caption:   media.Caption,
		SortOrder: media.SortOrder,
	}, nil
}

// DeleteProjectMedia deletes project media
func (s *ProjectService) DeleteProjectMedia(mediaID uint) error {
	var media models.ProjectMedia
	if err := database.DB.First(&media, mediaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("media not found")
		}
		return err
	}

	return database.DB.Delete(&media).Error
}

// Helper functions
func (s *ProjectService) mapProjectToResponse(project models.Project) *ProjectResponse {
	response := &ProjectResponse{
		ID:          project.ID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description,
		Content:     project.Content,
		CategoryID:  project.CategoryID,
		Featured:    project.Featured,
		Published:   project.Published,
		CreatedBy:   project.CreatedBy,
		UpdatedBy:   project.UpdatedBy,
		CreatedAt:   project.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   project.UpdatedAt.Format(time.RFC3339),
	}

	// Map category
	if project.Category.ID > 0 {
		response.Category = ProjectCategoryResponse{
			ID:   project.Category.ID,
			Name: project.Category.Name,
			Slug: project.Category.Slug,
		}
	}

	// Map media
	for _, media := range project.Media {
		response.Media = append(response.Media, ProjectMediaResponse{
			ID:        media.ID,
			Type:      media.Type,
			URL:       media.URL,
			Caption:   media.Caption,
			SortOrder: media.SortOrder,
		})
	}

	// Map tags
	for _, tag := range project.Tags {
		response.Tags = append(response.Tags, TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response
} 
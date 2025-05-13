package services

import (
	"errors"
	"strings"
	"time"

	"zionechainapi/internal/database"
	"zionechainapi/internal/models"
	"gorm.io/gorm"
)

// BlogService handles blog-related operations
type BlogService struct{}

// NewBlogService creates a new blog service
func NewBlogService() *BlogService {
	return &BlogService{}
}

// CreateBlogRequest represents the create blog request
type CreateBlogRequest struct {
	Title      string   `json:"title" binding:"required"`
	Excerpt    string   `json:"excerpt" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	CategoryID uint     `json:"category_id" binding:"required"`
	TagIDs     []uint   `json:"tag_ids"`
	Featured   bool     `json:"featured"`
	Published  bool     `json:"published"`
}

// UpdateBlogRequest represents the update blog request
type UpdateBlogRequest struct {
	Title      string  `json:"title"`
	Excerpt    string  `json:"excerpt"`
	Content    string  `json:"content"`
	CategoryID uint    `json:"category_id"`
	TagIDs     []uint  `json:"tag_ids"`
	Featured   *bool   `json:"featured"`
	Published  *bool   `json:"published"`
}

// BlogMediaRequest represents the blog media request
type BlogMediaRequest struct {
	Type      string `json:"type" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

// BlogResponse represents the blog response
type BlogResponse struct {
	ID         uint                 `json:"id"`
	Title      string               `json:"title"`
	Slug       string               `json:"slug"`
	Excerpt    string               `json:"excerpt"`
	Content    string               `json:"content"`
	CategoryID uint                 `json:"category_id"`
	Category   BlogCategoryResponse `json:"category"`
	Media      []BlogMediaResponse  `json:"media"`
	Tags       []TagResponse        `json:"tags"`
	Featured   bool                 `json:"featured"`
	Published  bool                 `json:"published"`
	CreatedBy  uint                 `json:"created_by"`
	UpdatedBy  uint                 `json:"updated_by"`
	CreatedAt  string               `json:"created_at"`
	UpdatedAt  string               `json:"updated_at"`
}

// BlogCategoryResponse represents the blog category response
type BlogCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// BlogMediaResponse represents the blog media response
type BlogMediaResponse struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

// CreateBlog creates a new blog post
func (s *BlogService) CreateBlog(req CreateBlogRequest, userID uint) (*BlogResponse, error) {
	// Create slug from title
	slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))

	// Check if slug already exists
	var count int64
	if err := database.DB.Model(&models.BlogPost{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		// Append timestamp to slug to make it unique
		slug = slug + "-" + string(time.Now().Unix())
	}

	// Create blog post
	blog := models.BlogPost{
		Title:      req.Title,
		Slug:       slug,
		Excerpt:    req.Excerpt,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		Featured:   req.Featured,
		Published:  req.Published,
		CreatedBy:  userID,
		UpdatedBy:  userID,
	}

	// Start transaction
	tx := database.DB.Begin()
	if err := tx.Create(&blog).Error; err != nil {
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

		if err := tx.Model(&blog).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load blog with relationships
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").First(&blog, blog.ID).Error; err != nil {
		return nil, err
	}

	// Map to response
	return s.mapBlogToResponse(blog), nil
}

// GetBlogByID gets a blog post by ID
func (s *BlogService) GetBlogByID(id uint) (*BlogResponse, error) {
	var blog models.BlogPost
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	return s.mapBlogToResponse(blog), nil
}

// GetBlogBySlug gets a blog post by slug
func (s *BlogService) GetBlogBySlug(slug string) (*BlogResponse, error) {
	var blog models.BlogPost
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").Where("slug = ?", slug).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	return s.mapBlogToResponse(blog), nil
}

// ListBlogs lists all blog posts with pagination
func (s *BlogService) ListBlogs(page, limit int, categoryID uint, featured, published bool) ([]BlogResponse, int64, error) {
	var blogs []models.BlogPost
	var total int64

	// Base query
	query := database.DB.Model(&models.BlogPost{})

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
		Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	// Map to response
	var response []BlogResponse
	for _, blog := range blogs {
		response = append(response, *s.mapBlogToResponse(blog))
	}

	return response, total, nil
}

// UpdateBlog updates a blog post
func (s *BlogService) UpdateBlog(id uint, req UpdateBlogRequest, userID uint) (*BlogResponse, error) {
	var blog models.BlogPost
	if err := database.DB.First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	// Update fields if provided
	tx := database.DB.Begin()

	if req.Title != "" && req.Title != blog.Title {
		// Create new slug from title
		slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))

		// Check if slug already exists and is not this blog
		var count int64
		if err := tx.Model(&models.BlogPost{}).Where("slug = ? AND id != ?", slug, id).Count(&count).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if count > 0 {
			// Append timestamp to slug to make it unique
			slug = slug + "-" + string(time.Now().Unix())
		}

		blog.Title = req.Title
		blog.Slug = slug
	}

	if req.Excerpt != "" {
		blog.Excerpt = req.Excerpt
	}

	if req.Content != "" {
		blog.Content = req.Content
	}

	if req.CategoryID > 0 {
		blog.CategoryID = req.CategoryID
	}

	if req.Featured != nil {
		blog.Featured = *req.Featured
	}

	if req.Published != nil {
		blog.Published = *req.Published
	}

	blog.UpdatedBy = userID

	if err := tx.Save(&blog).Error; err != nil {
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

		if err := tx.Model(&blog).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load blog with relationships
	if err := database.DB.Preload("Category").Preload("Media").Preload("Tags").First(&blog, id).Error; err != nil {
		return nil, err
	}

	return s.mapBlogToResponse(blog), nil
}

// DeleteBlog deletes a blog post
func (s *BlogService) DeleteBlog(id uint) error {
	var blog models.BlogPost
	if err := database.DB.First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("blog post not found")
		}
		return err
	}

	// Start transaction
	tx := database.DB.Begin()

	// Delete media
	if err := tx.Where("blog_id = ?", id).Delete(&models.BlogMedia{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Remove tag associations
	if err := tx.Model(&blog).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Delete blog
	if err := tx.Delete(&blog).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// AddBlogMedia adds media to a blog post
func (s *BlogService) AddBlogMedia(blogID uint, req BlogMediaRequest) (*BlogMediaResponse, error) {
	var blog models.BlogPost
	if err := database.DB.First(&blog, blogID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog post not found")
		}
		return nil, err
	}

	media := models.BlogMedia{
		BlogID:    blogID,
		Type:      req.Type,
		URL:       req.URL,
		Caption:   req.Caption,
		SortOrder: req.SortOrder,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		return nil, err
	}

	return &BlogMediaResponse{
		ID:        media.ID,
		Type:      media.Type,
		URL:       media.URL,
		Caption:   media.Caption,
		SortOrder: media.SortOrder,
	}, nil
}

// UpdateBlogMedia updates blog media
func (s *BlogService) UpdateBlogMedia(mediaID uint, req BlogMediaRequest) (*BlogMediaResponse, error) {
	var media models.BlogMedia
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

	return &BlogMediaResponse{
		ID:        media.ID,
		Type:      media.Type,
		URL:       media.URL,
		Caption:   media.Caption,
		SortOrder: media.SortOrder,
	}, nil
}

// DeleteBlogMedia deletes blog media
func (s *BlogService) DeleteBlogMedia(mediaID uint) error {
	var media models.BlogMedia
	if err := database.DB.First(&media, mediaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("media not found")
		}
		return err
	}

	return database.DB.Delete(&media).Error
}

// Helper functions
func (s *BlogService) mapBlogToResponse(blog models.BlogPost) *BlogResponse {
	response := &BlogResponse{
		ID:         blog.ID,
		Title:      blog.Title,
		Slug:       blog.Slug,
		Excerpt:    blog.Excerpt,
		Content:    blog.Content,
		CategoryID: blog.CategoryID,
		Featured:   blog.Featured,
		Published:  blog.Published,
		CreatedBy:  blog.CreatedBy,
		UpdatedBy:  blog.UpdatedBy,
		CreatedAt:  blog.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  blog.UpdatedAt.Format(time.RFC3339),
	}

	// Map category
	if blog.Category.ID > 0 {
		response.Category = BlogCategoryResponse{
			ID:   blog.Category.ID,
			Name: blog.Category.Name,
			Slug: blog.Category.Slug,
		}
	}

	// Map media
	for _, media := range blog.Media {
		response.Media = append(response.Media, BlogMediaResponse{
			ID:        media.ID,
			Type:      media.Type,
			URL:       media.URL,
			Caption:   media.Caption,
			SortOrder: media.SortOrder,
		})
	}

	// Map tags
	for _, tag := range blog.Tags {
		response.Tags = append(response.Tags, TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return response
} 
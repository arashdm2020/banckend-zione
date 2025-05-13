package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/middleware"
	"zionechainapi/internal/services"
	"zionechainapi/internal/utils"
)

// BlogController handles blog-related routes
type BlogController struct {
	config      *configs.Config
	blogService *services.BlogService
}

// NewBlogController creates a new blog controller
func NewBlogController(config *configs.Config) *BlogController {
	return &BlogController{
		config:      config,
		blogService: services.NewBlogService(),
	}
}

// Create godoc
// @Summary Create a new blog post
// @Description Create a new blog post
// @Tags blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.CreateBlogRequest true "Create blog request"
// @Success 201 {object} utils.Response{data=services.BlogResponse} "Blog post created successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog [post]
func (c *BlogController) Create(ctx *gin.Context) {
	var req services.CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	userID := middleware.GetUserID(ctx)
	blog, err := c.blogService.CreateBlog(req, userID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to create blog post", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "Blog post created successfully", blog)
}

// Get godoc
// @Summary Get a blog post by ID
// @Description Get a blog post by ID
// @Tags blog
// @Accept json
// @Produce json
// @Param id path int true "Blog Post ID"
// @Success 200 {object} utils.Response{data=services.BlogResponse} "Blog post retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/{id} [get]
func (c *BlogController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid blog post ID", nil)
		return
	}

	blog, err := c.blogService.GetBlogByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Blog post retrieved successfully", blog)
}

// GetBySlug godoc
// @Summary Get a blog post by slug
// @Description Get a blog post by slug
// @Tags blog
// @Accept json
// @Produce json
// @Param slug path string true "Blog Post Slug"
// @Success 200 {object} utils.Response{data=services.BlogResponse} "Blog post retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/slug/{slug} [get]
func (c *BlogController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	blog, err := c.blogService.GetBlogBySlug(slug)
	if err != nil {
		utils.NotFoundResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Blog post retrieved successfully", blog)
}

// List godoc
// @Summary List blog posts
// @Description List blog posts with pagination
// @Tags blog
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param category_id query int false "Category ID"
// @Param featured query bool false "Featured flag"
// @Success 200 {object} utils.Response{data=[]services.BlogResponse} "Blog posts retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog [get]
func (c *BlogController) List(ctx *gin.Context) {
	page := 1
	limit := 10
	var categoryID uint
	featured := false
	published := true // Default to published only

	// Parse query parameters
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil {
			page = pageNum
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil {
			limit = limitNum
		}
	}

	if categoryIDStr := ctx.Query("category_id"); categoryIDStr != "" {
		if categoryIDNum, err := strconv.ParseUint(categoryIDStr, 10, 64); err == nil {
			categoryID = uint(categoryIDNum)
		}
	}

	if featuredStr := ctx.Query("featured"); featuredStr != "" {
		if featuredBool, err := strconv.ParseBool(featuredStr); err == nil {
			featured = featuredBool
		}
	}

	// Check if user is admin or editor
	userRole := middleware.GetUserRole(ctx)
	if userRole == "admin" || userRole == "editor" {
		// If user is admin or editor, check if they want to see unpublished blog posts
		if publishedStr := ctx.Query("published"); publishedStr != "" {
			if publishedBool, err := strconv.ParseBool(publishedStr); err == nil {
				published = publishedBool
			}
		}
	}

	blogs, total, err := c.blogService.ListBlogs(page, limit, categoryID, featured, published)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	// Create response with pagination metadata
	response := map[string]interface{}{
		"blogs": blogs,
		"metadata": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.OKResponse(ctx, "Blog posts retrieved successfully", response)
}

// Update godoc
// @Summary Update a blog post
// @Description Update a blog post
// @Tags blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Blog Post ID"
// @Param body body services.UpdateBlogRequest true "Update blog request"
// @Success 200 {object} utils.Response{data=services.BlogResponse} "Blog post updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/{id} [put]
func (c *BlogController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid blog post ID", nil)
		return
	}

	var req services.UpdateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	userID := middleware.GetUserID(ctx)
	blog, err := c.blogService.UpdateBlog(uint(id), req, userID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to update blog post", err.Error())
		return
	}

	utils.OKResponse(ctx, "Blog post updated successfully", blog)
}

// Delete godoc
// @Summary Delete a blog post
// @Description Delete a blog post
// @Tags blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Blog Post ID"
// @Success 204 {object} utils.Response "Blog post deleted successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/{id} [delete]
func (c *BlogController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid blog post ID", nil)
		return
	}

	if err := c.blogService.DeleteBlog(uint(id)); err != nil {
		utils.BadRequestResponse(ctx, "Failed to delete blog post", err.Error())
		return
	}

	utils.NoContentResponse(ctx)
}

// AddMedia godoc
// @Summary Add media to a blog post
// @Description Add media to a blog post
// @Tags blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Blog Post ID"
// @Param body body services.BlogMediaRequest true "Add media request"
// @Success 201 {object} utils.Response{data=services.BlogMediaResponse} "Media added successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/{id}/media [post]
func (c *BlogController) AddMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid blog post ID", nil)
		return
	}

	var req services.BlogMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	media, err := c.blogService.AddBlogMedia(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to add media", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "Media added successfully", media)
}

// UpdateMedia godoc
// @Summary Update blog media
// @Description Update blog media
// @Tags blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Media ID"
// @Param body body services.BlogMediaRequest true "Update media request"
// @Success 200 {object} utils.Response{data=services.BlogMediaResponse} "Media updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/media/{id} [put]
func (c *BlogController) UpdateMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid media ID", nil)
		return
	}

	var req services.BlogMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	media, err := c.blogService.UpdateBlogMedia(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to update media", err.Error())
		return
	}

	utils.OKResponse(ctx, "Media updated successfully", media)
}

// DeleteMedia godoc
// @Summary Delete blog media
// @Description Delete blog media
// @Tags blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Media ID"
// @Success 204 {object} utils.Response "Media deleted successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/blog/media/{id} [delete]
func (c *BlogController) DeleteMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid media ID", nil)
		return
	}

	if err := c.blogService.DeleteBlogMedia(uint(id)); err != nil {
		utils.BadRequestResponse(ctx, "Failed to delete media", err.Error())
		return
	}

	utils.NoContentResponse(ctx)
}

// Routes registers blog routes
func (c *BlogController) Routes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	blog := router.Group("/blog")
	{
		// Public routes
		blog.GET("", c.List)
		blog.GET("/:id", c.Get)
		blog.GET("/slug/:slug", c.GetBySlug)

		// Protected routes
		authenticated := blog.Group("")
		authenticated.Use(authMiddleware)
		{
			// Admin and editor routes
			adminEditor := authenticated.Group("")
			adminEditor.Use(middleware.RequireRole("admin", "editor"))
			{
				adminEditor.POST("", c.Create)
				adminEditor.PUT("/:id", c.Update)
				adminEditor.DELETE("/:id", c.Delete)
				adminEditor.POST("/:id/media", c.AddMedia)
				adminEditor.PUT("/media/:id", c.UpdateMedia)
				adminEditor.DELETE("/media/:id", c.DeleteMedia)
			}
		}
	}
} 
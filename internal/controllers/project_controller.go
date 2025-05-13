package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/middleware"
	"zionechainapi/internal/services"
	"zionechainapi/internal/utils"
)

// ProjectController handles project-related routes
type ProjectController struct {
	config         *configs.Config
	projectService *services.ProjectService
}

// NewProjectController creates a new project controller
func NewProjectController(config *configs.Config) *ProjectController {
	return &ProjectController{
		config:         config,
		projectService: services.NewProjectService(),
	}
}

// Create godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.CreateProjectRequest true "Create project request"
// @Success 201 {object} utils.Response{data=services.ProjectResponse} "Project created successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects [post]
func (c *ProjectController) Create(ctx *gin.Context) {
	var req services.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	userID := middleware.GetUserID(ctx)
	project, err := c.projectService.CreateProject(req, userID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to create project", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "Project created successfully", project)
}

// Get godoc
// @Summary Get a project by ID
// @Description Get a project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} utils.Response{data=services.ProjectResponse} "Project retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects/{id} [get]
func (c *ProjectController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid project ID", nil)
		return
	}

	project, err := c.projectService.GetProjectByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Project retrieved successfully", project)
}

// GetBySlug godoc
// @Summary Get a project by slug
// @Description Get a project by slug
// @Tags projects
// @Accept json
// @Produce json
// @Param slug path string true "Project Slug"
// @Success 200 {object} utils.Response{data=services.ProjectResponse} "Project retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects/slug/{slug} [get]
func (c *ProjectController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	project, err := c.projectService.GetProjectBySlug(slug)
	if err != nil {
		utils.NotFoundResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Project retrieved successfully", project)
}

// List godoc
// @Summary List projects
// @Description List projects with pagination
// @Tags projects
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param category_id query int false "Category ID"
// @Param featured query bool false "Featured flag"
// @Success 200 {object} utils.Response{data=[]services.ProjectResponse} "Projects retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects [get]
func (c *ProjectController) List(ctx *gin.Context) {
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
		// If user is admin or editor, check if they want to see unpublished projects
		if publishedStr := ctx.Query("published"); publishedStr != "" {
			if publishedBool, err := strconv.ParseBool(publishedStr); err == nil {
				published = publishedBool
			}
		}
	}

	projects, total, err := c.projectService.ListProjects(page, limit, categoryID, featured, published)
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	// Create response with pagination metadata
	response := map[string]interface{}{
		"projects": projects,
		"metadata": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.OKResponse(ctx, "Projects retrieved successfully", response)
}

// Update godoc
// @Summary Update a project
// @Description Update a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param body body services.UpdateProjectRequest true "Update project request"
// @Success 200 {object} utils.Response{data=services.ProjectResponse} "Project updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects/{id} [put]
func (c *ProjectController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid project ID", nil)
		return
	}

	var req services.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	userID := middleware.GetUserID(ctx)
	project, err := c.projectService.UpdateProject(uint(id), req, userID)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to update project", err.Error())
		return
	}

	utils.OKResponse(ctx, "Project updated successfully", project)
}

// Delete godoc
// @Summary Delete a project
// @Description Delete a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 204 {object} utils.Response "Project deleted successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects/{id} [delete]
func (c *ProjectController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid project ID", nil)
		return
	}

	if err := c.projectService.DeleteProject(uint(id)); err != nil {
		utils.BadRequestResponse(ctx, "Failed to delete project", err.Error())
		return
	}

	utils.NoContentResponse(ctx)
}

// AddMedia godoc
// @Summary Add media to a project
// @Description Add media to a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param body body services.ProjectMediaRequest true "Add media request"
// @Success 201 {object} utils.Response{data=services.ProjectMediaResponse} "Media added successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects/{id}/media [post]
func (c *ProjectController) AddMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid project ID", nil)
		return
	}

	var req services.ProjectMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	media, err := c.projectService.AddProjectMedia(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to add media", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "Media added successfully", media)
}

// UpdateMedia godoc
// @Summary Update project media
// @Description Update project media
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Media ID"
// @Param body body services.ProjectMediaRequest true "Update media request"
// @Success 200 {object} utils.Response{data=services.ProjectMediaResponse} "Media updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/projects/media/{id} [put]
func (c *ProjectController) UpdateMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid media ID", nil)
		return
	}

	var req services.ProjectMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	media, err := c.projectService.UpdateProjectMedia(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to update media", err.Error())
		return
	}

	utils.OKResponse(ctx, "Media updated successfully", media)
}

// DeleteMedia godoc
// @Summary Delete project media
// @Description Delete project media
// @Tags projects
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
// @Router /api/projects/media/{id} [delete]
func (c *ProjectController) DeleteMedia(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid media ID", nil)
		return
	}

	if err := c.projectService.DeleteProjectMedia(uint(id)); err != nil {
		utils.BadRequestResponse(ctx, "Failed to delete media", err.Error())
		return
	}

	utils.NoContentResponse(ctx)
}

// Routes registers project routes
func (c *ProjectController) Routes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	projects := router.Group("/projects")
	{
		// Public routes
		projects.GET("", c.List)
		projects.GET("/:id", c.Get)
		projects.GET("/slug/:slug", c.GetBySlug)

		// Protected routes
		authenticated := projects.Group("")
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
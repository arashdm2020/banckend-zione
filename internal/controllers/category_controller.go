package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"zionechainapi/configs"
	"zionechainapi/internal/middleware"
	"zionechainapi/internal/services"
	"zionechainapi/internal/utils"
)

// CategoryController handles category-related routes
type CategoryController struct {
	config          *configs.Config
	categoryService *services.CategoryService
}

// NewCategoryController creates a new category controller
func NewCategoryController(config *configs.Config) *CategoryController {
	return &CategoryController{
		config:          config,
		categoryService: services.NewCategoryService(),
	}
}

// CreateProjectCategory godoc
// @Summary Create a new project category
// @Description Create a new project category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.CategoryRequest true "Create category request"
// @Success 201 {object} utils.Response{data=services.ProjectCategoryResponse} "Category created successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/projects [post]
func (c *CategoryController) CreateProjectCategory(ctx *gin.Context) {
	var req services.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	category, err := c.categoryService.CreateProjectCategory(req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to create project category", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "Project category created successfully", category)
}

// ListProjectCategories godoc
// @Summary List all project categories
// @Description List all project categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]services.ProjectCategoryResponse} "Categories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/projects [get]
func (c *CategoryController) ListProjectCategories(ctx *gin.Context) {
	categories, err := c.categoryService.ListProjectCategories()
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Project categories retrieved successfully", categories)
}

// GetProjectCategory godoc
// @Summary Get a project category by ID
// @Description Get a project category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response{data=services.ProjectCategoryResponse} "Category retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/projects/{id} [get]
func (c *CategoryController) GetProjectCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", nil)
		return
	}

	category, err := c.categoryService.GetProjectCategoryByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Project category retrieved successfully", category)
}

// UpdateProjectCategory godoc
// @Summary Update a project category
// @Description Update a project category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param body body services.CategoryRequest true "Update category request"
// @Success 200 {object} utils.Response{data=services.ProjectCategoryResponse} "Category updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/projects/{id} [put]
func (c *CategoryController) UpdateProjectCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", nil)
		return
	}

	var req services.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	category, err := c.categoryService.UpdateProjectCategory(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to update project category", err.Error())
		return
	}

	utils.OKResponse(ctx, "Project category updated successfully", category)
}

// DeleteProjectCategory godoc
// @Summary Delete a project category
// @Description Delete a project category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204 {object} utils.Response "Category deleted successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/projects/{id} [delete]
func (c *CategoryController) DeleteProjectCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", nil)
		return
	}

	if err := c.categoryService.DeleteProjectCategory(uint(id)); err != nil {
		utils.BadRequestResponse(ctx, "Failed to delete project category", err.Error())
		return
	}

	utils.NoContentResponse(ctx)
}

// CreateBlogCategory godoc
// @Summary Create a new blog category
// @Description Create a new blog category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body services.CategoryRequest true "Create category request"
// @Success 201 {object} utils.Response{data=services.BlogCategoryResponse} "Category created successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/blog [post]
func (c *CategoryController) CreateBlogCategory(ctx *gin.Context) {
	var req services.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	category, err := c.categoryService.CreateBlogCategory(req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to create blog category", err.Error())
		return
	}

	utils.CreatedResponse(ctx, "Blog category created successfully", category)
}

// ListBlogCategories godoc
// @Summary List all blog categories
// @Description List all blog categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]services.BlogCategoryResponse} "Categories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/blog [get]
func (c *CategoryController) ListBlogCategories(ctx *gin.Context) {
	categories, err := c.categoryService.ListBlogCategories()
	if err != nil {
		utils.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Blog categories retrieved successfully", categories)
}

// GetBlogCategory godoc
// @Summary Get a blog category by ID
// @Description Get a blog category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response{data=services.BlogCategoryResponse} "Category retrieved successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/blog/{id} [get]
func (c *CategoryController) GetBlogCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", nil)
		return
	}

	category, err := c.categoryService.GetBlogCategoryByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(ctx, err.Error())
		return
	}

	utils.OKResponse(ctx, "Blog category retrieved successfully", category)
}

// UpdateBlogCategory godoc
// @Summary Update a blog category
// @Description Update a blog category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param body body services.CategoryRequest true "Update category request"
// @Success 200 {object} utils.Response{data=services.BlogCategoryResponse} "Category updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 422 {object} utils.Response "Validation error"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/blog/{id} [put]
func (c *CategoryController) UpdateBlogCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", nil)
		return
	}

	var req services.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	category, err := c.categoryService.UpdateBlogCategory(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(ctx, "Failed to update blog category", err.Error())
		return
	}

	utils.OKResponse(ctx, "Blog category updated successfully", category)
}

// DeleteBlogCategory godoc
// @Summary Delete a blog category
// @Description Delete a blog category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204 {object} utils.Response "Category deleted successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
// @Failure 404 {object} utils.Response "Not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /api/categories/blog/{id} [delete]
func (c *CategoryController) DeleteBlogCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(ctx, "Invalid category ID", nil)
		return
	}

	if err := c.categoryService.DeleteBlogCategory(uint(id)); err != nil {
		utils.BadRequestResponse(ctx, "Failed to delete blog category", err.Error())
		return
	}

	utils.NoContentResponse(ctx)
}

// Routes registers category routes
func (c *CategoryController) Routes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	categories := router.Group("/categories")
	{
		// Project category routes
		projectCategories := categories.Group("/projects")
		{
			// Public routes
			projectCategories.GET("", c.ListProjectCategories)
			projectCategories.GET("/:id", c.GetProjectCategory)

			// Protected routes
			authenticated := projectCategories.Group("")
			authenticated.Use(authMiddleware)
			{
				// Admin and editor routes
				adminEditor := authenticated.Group("")
				adminEditor.Use(middleware.RequireRole("admin", "editor"))
				{
					adminEditor.POST("", c.CreateProjectCategory)
					adminEditor.PUT("/:id", c.UpdateProjectCategory)
					adminEditor.DELETE("/:id", c.DeleteProjectCategory)
				}
			}
		}

		// Blog category routes
		blogCategories := categories.Group("/blog")
		{
			// Public routes
			blogCategories.GET("", c.ListBlogCategories)
			blogCategories.GET("/:id", c.GetBlogCategory)

			// Protected routes
			authenticated := blogCategories.Group("")
			authenticated.Use(authMiddleware)
			{
				// Admin and editor routes
				adminEditor := authenticated.Group("")
				adminEditor.Use(middleware.RequireRole("admin", "editor"))
				{
					adminEditor.POST("", c.CreateBlogCategory)
					adminEditor.PUT("/:id", c.UpdateBlogCategory)
					adminEditor.DELETE("/:id", c.DeleteBlogCategory)
				}
			}
		}
	}
} 
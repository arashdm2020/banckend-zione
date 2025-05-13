package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/arashdm2020/banckend-zione/internal/models"
)

// ResumeController handles resume-related API requests
type ResumeController struct {
	DB *gorm.DB
}

// NewResumeController creates a new resume controller
func NewResumeController(db *gorm.DB) *ResumeController {
	return &ResumeController{
		DB: db,
	}
}

// Routes sets up the resume routes
func (c *ResumeController) Routes(router *gin.RouterGroup) {
	resumeRoutes := router.Group("/resume")
	{
		// Personal Info
		resumeRoutes.GET("/personal", c.GetPersonalInfo)
		resumeRoutes.POST("/personal", c.CreatePersonalInfo)
		resumeRoutes.PUT("/personal/:id", c.UpdatePersonalInfo)
		resumeRoutes.DELETE("/personal/:id", c.DeletePersonalInfo)

		// Skills
		resumeRoutes.GET("/skills", c.GetSkills)
		resumeRoutes.POST("/skills", c.CreateSkill)
		resumeRoutes.PUT("/skills/:id", c.UpdateSkill)
		resumeRoutes.DELETE("/skills/:id", c.DeleteSkill)

		// Experience
		resumeRoutes.GET("/experience", c.GetExperiences)
		resumeRoutes.POST("/experience", c.CreateExperience)
		resumeRoutes.PUT("/experience/:id", c.UpdateExperience)
		resumeRoutes.DELETE("/experience/:id", c.DeleteExperience)

		// Education
		resumeRoutes.GET("/education", c.GetEducations)
		resumeRoutes.POST("/education", c.CreateEducation)
		resumeRoutes.PUT("/education/:id", c.UpdateEducation)
		resumeRoutes.DELETE("/education/:id", c.DeleteEducation)

		// Projects
		resumeRoutes.GET("/projects", c.GetProjects)
		resumeRoutes.POST("/projects", c.CreateProject)
		resumeRoutes.PUT("/projects/:id", c.UpdateProject)
		resumeRoutes.DELETE("/projects/:id", c.DeleteProject)

		// Certificates
		resumeRoutes.GET("/certificates", c.GetCertificates)
		resumeRoutes.POST("/certificates", c.CreateCertificate)
		resumeRoutes.PUT("/certificates/:id", c.UpdateCertificate)
		resumeRoutes.DELETE("/certificates/:id", c.DeleteCertificate)

		// Languages
		resumeRoutes.GET("/languages", c.GetLanguages)
		resumeRoutes.POST("/languages", c.CreateLanguage)
		resumeRoutes.PUT("/languages/:id", c.UpdateLanguage)
		resumeRoutes.DELETE("/languages/:id", c.DeleteLanguage)

		// Publications
		resumeRoutes.GET("/publications", c.GetPublications)
		resumeRoutes.POST("/publications", c.CreatePublication)
		resumeRoutes.PUT("/publications/:id", c.UpdatePublication)
		resumeRoutes.DELETE("/publications/:id", c.DeletePublication)

		// Complete Resume
		resumeRoutes.GET("/complete", c.GetCompleteResume)
	}
}

// GetCompleteResume returns all resume sections
func (c *ResumeController) GetCompleteResume(ctx *gin.Context) {
	var personalInfo []models.PersonalInfo
	var skills []models.Skill
	var experiences []models.Experience
	var educations []models.Education
	var projects []models.Project
	var certificates []models.Certificate
	var languages []models.Language
	var publications []models.Publication

	c.DB.Find(&personalInfo)
	c.DB.Find(&skills)
	c.DB.Find(&experiences)
	c.DB.Find(&educations)
	c.DB.Find(&projects)
	c.DB.Find(&certificates)
	c.DB.Find(&languages)
	c.DB.Find(&publications)

	response := gin.H{
		"personal_info": personalInfo,
		"skills":        skills,
		"experience":    experiences,
		"education":     educations,
		"projects":      projects,
		"certificates":  certificates,
		"languages":     languages,
		"publications":  publications,
	}

	ctx.JSON(http.StatusOK, response)
}

// Personal Info controller methods
func (c *ResumeController) GetPersonalInfo(ctx *gin.Context) {
	var personalInfo []models.PersonalInfo
	c.DB.Find(&personalInfo)
	ctx.JSON(http.StatusOK, personalInfo)
}

func (c *ResumeController) CreatePersonalInfo(ctx *gin.Context) {
	var input models.PersonalInfo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdatePersonalInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	var personalInfo models.PersonalInfo
	if err := c.DB.First(&personalInfo, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.PersonalInfo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&personalInfo).Updates(input)
	ctx.JSON(http.StatusOK, personalInfo)
}

func (c *ResumeController) DeletePersonalInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	var personalInfo models.PersonalInfo
	if err := c.DB.First(&personalInfo, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&personalInfo)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Skills controller methods
func (c *ResumeController) GetSkills(ctx *gin.Context) {
	var skills []models.Skill
	c.DB.Find(&skills)
	ctx.JSON(http.StatusOK, skills)
}

func (c *ResumeController) CreateSkill(ctx *gin.Context) {
	var input models.Skill
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdateSkill(ctx *gin.Context) {
	id := ctx.Param("id")
	var skill models.Skill
	if err := c.DB.First(&skill, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Skill
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&skill).Updates(input)
	ctx.JSON(http.StatusOK, skill)
}

func (c *ResumeController) DeleteSkill(ctx *gin.Context) {
	id := ctx.Param("id")
	var skill models.Skill
	if err := c.DB.First(&skill, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&skill)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Experience controller methods
func (c *ResumeController) GetExperiences(ctx *gin.Context) {
	var experiences []models.Experience
	c.DB.Find(&experiences)
	ctx.JSON(http.StatusOK, experiences)
}

func (c *ResumeController) CreateExperience(ctx *gin.Context) {
	var input models.Experience
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdateExperience(ctx *gin.Context) {
	id := ctx.Param("id")
	var experience models.Experience
	if err := c.DB.First(&experience, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Experience
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&experience).Updates(input)
	ctx.JSON(http.StatusOK, experience)
}

func (c *ResumeController) DeleteExperience(ctx *gin.Context) {
	id := ctx.Param("id")
	var experience models.Experience
	if err := c.DB.First(&experience, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&experience)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Education controller methods
func (c *ResumeController) GetEducations(ctx *gin.Context) {
	var educations []models.Education
	c.DB.Find(&educations)
	ctx.JSON(http.StatusOK, educations)
}

func (c *ResumeController) CreateEducation(ctx *gin.Context) {
	var input models.Education
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdateEducation(ctx *gin.Context) {
	id := ctx.Param("id")
	var education models.Education
	if err := c.DB.First(&education, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Education
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&education).Updates(input)
	ctx.JSON(http.StatusOK, education)
}

func (c *ResumeController) DeleteEducation(ctx *gin.Context) {
	id := ctx.Param("id")
	var education models.Education
	if err := c.DB.First(&education, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&education)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Project controller methods
func (c *ResumeController) GetProjects(ctx *gin.Context) {
	var projects []models.Project
	c.DB.Find(&projects)
	ctx.JSON(http.StatusOK, projects)
}

func (c *ResumeController) CreateProject(ctx *gin.Context) {
	var input models.Project
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdateProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project
	if err := c.DB.First(&project, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Project
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&project).Updates(input)
	ctx.JSON(http.StatusOK, project)
}

func (c *ResumeController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project
	if err := c.DB.First(&project, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&project)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Certificate controller methods
func (c *ResumeController) GetCertificates(ctx *gin.Context) {
	var certificates []models.Certificate
	c.DB.Find(&certificates)
	ctx.JSON(http.StatusOK, certificates)
}

func (c *ResumeController) CreateCertificate(ctx *gin.Context) {
	var input models.Certificate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdateCertificate(ctx *gin.Context) {
	id := ctx.Param("id")
	var certificate models.Certificate
	if err := c.DB.First(&certificate, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Certificate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&certificate).Updates(input)
	ctx.JSON(http.StatusOK, certificate)
}

func (c *ResumeController) DeleteCertificate(ctx *gin.Context) {
	id := ctx.Param("id")
	var certificate models.Certificate
	if err := c.DB.First(&certificate, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&certificate)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Language controller methods
func (c *ResumeController) GetLanguages(ctx *gin.Context) {
	var languages []models.Language
	c.DB.Find(&languages)
	ctx.JSON(http.StatusOK, languages)
}

func (c *ResumeController) CreateLanguage(ctx *gin.Context) {
	var input models.Language
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdateLanguage(ctx *gin.Context) {
	id := ctx.Param("id")
	var language models.Language
	if err := c.DB.First(&language, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Language
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&language).Updates(input)
	ctx.JSON(http.StatusOK, language)
}

func (c *ResumeController) DeleteLanguage(ctx *gin.Context) {
	id := ctx.Param("id")
	var language models.Language
	if err := c.DB.First(&language, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&language)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// Publication controller methods
func (c *ResumeController) GetPublications(ctx *gin.Context) {
	var publications []models.Publication
	c.DB.Find(&publications)
	ctx.JSON(http.StatusOK, publications)
}

func (c *ResumeController) CreatePublication(ctx *gin.Context) {
	var input models.Publication
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Create(&input)
	ctx.JSON(http.StatusCreated, input)
}

func (c *ResumeController) UpdatePublication(ctx *gin.Context) {
	id := ctx.Param("id")
	var publication models.Publication
	if err := c.DB.First(&publication, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var input models.Publication
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.DB.Model(&publication).Updates(input)
	ctx.JSON(http.StatusOK, publication)
}

func (c *ResumeController) DeletePublication(ctx *gin.Context) {
	id := ctx.Param("id")
	var publication models.Publication
	if err := c.DB.First(&publication, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.DB.Delete(&publication)
	ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
} 
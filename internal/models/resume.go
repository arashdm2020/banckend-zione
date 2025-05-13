package models

import (
	"time"

	"gorm.io/gorm"
)

// PersonalInfo represents personal information section in resume
type PersonalInfo struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	FullName     string         `json:"full_name" binding:"required"`
	JobTitle     string         `json:"job_title" binding:"required"`
	Email        string         `json:"email" binding:"required,email"`
	Phone        string         `json:"phone" binding:"required"`
	Address      string         `json:"address"`
	Website      string         `json:"website"`
	LinkedIn     string         `json:"linkedin"`
	GitHub       string         `json:"github"`
	Twitter      string         `json:"twitter"`
	Summary      string         `json:"summary" binding:"required"`
	ProfileImage string         `json:"profile_image"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Skill represents skill section in resume
type Skill struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" binding:"required"`
	Proficiency int            `json:"proficiency" binding:"required,min=1,max=100"`
	Category    string         `json:"category"`
	IconURL     string         `json:"icon_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Experience represents work experience section in resume
type Experience struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	JobTitle     string         `json:"job_title" binding:"required"`
	Company      string         `json:"company" binding:"required"`
	Location     string         `json:"location"`
	StartDate    time.Time      `json:"start_date" binding:"required"`
	EndDate      *time.Time     `json:"end_date"`
	CurrentJob   bool           `json:"current_job"`
	Description  string         `json:"description" binding:"required"`
	Achievements string         `json:"achievements"`
	Website      string         `json:"website"`
	LogoURL      string         `json:"logo_url"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Education represents education section in resume
type Education struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Institution string         `json:"institution" binding:"required"`
	Degree      string         `json:"degree" binding:"required"`
	Field       string         `json:"field" binding:"required"`
	Location    string         `json:"location"`
	StartDate   time.Time      `json:"start_date" binding:"required"`
	EndDate     *time.Time     `json:"end_date"`
	Current     bool           `json:"current"`
	GPA         string         `json:"gpa"`
	Description string         `json:"description"`
	LogoURL     string         `json:"logo_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Project represents project section in resume
type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Role        string         `json:"role"`
	URL         string         `json:"url"`
	GitHubURL   string         `json:"github_url"`
	ImageURL    string         `json:"image_url"`
	Technologies string        `json:"technologies"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Ongoing     bool           `json:"ongoing"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Certificate represents certification section in resume
type Certificate struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" binding:"required"`
	Issuer      string         `json:"issuer" binding:"required"`
	IssueDate   time.Time      `json:"issue_date" binding:"required"`
	ExpiryDate  *time.Time     `json:"expiry_date"`
	NoExpiry    bool           `json:"no_expiry"`
	CredentialID string        `json:"credential_id"`
	CredentialURL string       `json:"credential_url"`
	Description string         `json:"description"`
	LogoURL     string         `json:"logo_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Language represents language proficiency section in resume
type Language struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" binding:"required"`
	Proficiency string         `json:"proficiency" binding:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Publication represents publications section in resume
type Publication struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" binding:"required"`
	Publisher   string         `json:"publisher" binding:"required"`
	Authors     string         `json:"authors"`
	PublishDate time.Time      `json:"publish_date" binding:"required"`
	URL         string         `json:"url"`
	DOI         string         `json:"doi"`
	Description string         `json:"description"`
	ImageURL    string         `json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
} 
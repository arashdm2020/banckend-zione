package models

import "time"

// Project represents a project in the portfolio
type Project struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Title       string          `gorm:"size:200;not null" json:"title"`
	Slug        string          `gorm:"size:200;not null;uniqueIndex" json:"slug"`
	Description string          `gorm:"type:text" json:"description"`
	Content     string          `gorm:"type:longtext" json:"content"`
	CategoryID  uint            `json:"category_id"`
	Category    ProjectCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Media       []ProjectMedia  `json:"media"`
	Tags        []Tag           `gorm:"many2many:project_tags;" json:"tags"`
	Featured    bool            `gorm:"default:false" json:"featured"`
	Published   bool            `gorm:"default:true" json:"published"`
	CreatedBy   uint            `json:"created_by"`
	UpdatedBy   uint            `json:"updated_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// TableName specifies the table name for Project
func (Project) TableName() string {
	return "projects"
}

// ProjectCategory represents a category for projects
type ProjectCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Slug      string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Projects  []Project `gorm:"foreignKey:CategoryID" json:"-"` // Avoid circular reference in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for ProjectCategory
func (ProjectCategory) TableName() string {
	return "project_categories"
}

// ProjectMedia represents media attached to a project
type ProjectMedia struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `json:"project_id"`
	Type      string    `gorm:"size:20;not null;default:'image'" json:"type"` // image, video, etc.
	URL       string    `gorm:"size:255;not null" json:"url"`
	Caption   string    `gorm:"size:255" json:"caption"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for ProjectMedia
func (ProjectMedia) TableName() string {
	return "project_media"
} 
package models

import "time"

// Tag represents a tag that can be associated with projects or blog posts
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Slug      string    `gorm:"size:50;not null;uniqueIndex" json:"slug"`
	Projects  []Project `gorm:"many2many:project_tags;" json:"-"`
	BlogPosts []BlogPost `gorm:"many2many:blog_tags;" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Tag
func (Tag) TableName() string {
	return "tags"
} 
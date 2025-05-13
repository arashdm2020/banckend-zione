package models

import "time"

// BlogPost represents a blog post
type BlogPost struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Title      string       `gorm:"size:200;not null" json:"title"`
	Slug       string       `gorm:"size:200;not null;uniqueIndex" json:"slug"`
	Excerpt    string       `gorm:"type:text" json:"excerpt"`
	Content    string       `gorm:"type:longtext" json:"content"`
	CategoryID uint         `json:"category_id"`
	Category   BlogCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Media      []BlogMedia  `json:"media"`
	Tags       []Tag        `gorm:"many2many:blog_tags;" json:"tags"`
	Featured   bool         `gorm:"default:false" json:"featured"`
	Published  bool         `gorm:"default:true" json:"published"`
	CreatedBy  uint         `json:"created_by"`
	UpdatedBy  uint         `json:"updated_by"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// TableName specifies the table name for BlogPost
func (BlogPost) TableName() string {
	return "blog_posts"
}

// BlogCategory represents a category for blog posts
type BlogCategory struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	Slug      string     `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Posts     []BlogPost `gorm:"foreignKey:CategoryID" json:"-"` // Avoid circular reference in JSON
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName specifies the table name for BlogCategory
func (BlogCategory) TableName() string {
	return "blog_categories"
}

// BlogMedia represents media attached to a blog post
type BlogMedia struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BlogID    uint      `json:"blog_id"`
	Type      string    `gorm:"size:20;not null;default:'image'" json:"type"` // image, video, etc.
	URL       string    `gorm:"size:255;not null" json:"url"`
	Caption   string    `gorm:"size:255" json:"caption"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for BlogMedia
func (BlogMedia) TableName() string {
	return "blog_media"
} 
package models

import "time"

// Role represents a user role in the system
type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}

// Predefined role constants
const (
	RoleAdmin  = 1
	RoleEditor = 2
	RoleUser   = 3
) 
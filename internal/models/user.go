package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Phone     string    `gorm:"size:15;not null;uniqueIndex" json:"phone"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RoleID    uint      `gorm:"not null;default:3" json:"role_id"` // Default to user role (3)
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook is called before creating a User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	// Hash password if it's not already hashed
	if len(u.Password) > 0 && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return err
}

// BeforeUpdate hook is called before updating a User
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	var err error
	// Hash password if it's not already hashed and if it's being updated
	if tx.Statement.Changed("Password") && len(u.Password) > 0 && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return err
}

// CheckPassword checks if the provided password is correct
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin checks if the user has the admin role
func (u *User) IsAdmin() bool {
	return u.RoleID == 1
}

// IsEditor checks if the user has the editor role
func (u *User) IsEditor() bool {
	return u.RoleID == 2
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(roleName string) bool {
	switch roleName {
	case "admin":
		return u.RoleID == 1
	case "editor":
		return u.RoleID == 2 || u.RoleID == 1 // Admin also has editor privileges
	case "user":
		return true // All authenticated users have user privileges
	default:
		return false
	}
} 
package database

import (
	"context"
	"fmt"
	"time"

	"zionechainapi/configs"
	"zionechainapi/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection
var DB *gorm.DB

// Connect connects to the database
func Connect(config *configs.Config) (*gorm.DB, error) {
	// Create DSN string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Charset,
	)

	// Set up GORM logger
	var logLevel logger.LogLevel
	switch config.Log.Level {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Info
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// Connect to database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.Database.ConnMaxLifetime)

	return DB, nil
}

// AutoMigrate automatically migrates the database schema
func AutoMigrate() error {
	// Register models here
	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Project{},
		&models.ProjectCategory{},
		&models.ProjectMedia{},
		&models.BlogPost{},
		&models.BlogCategory{},
		&models.BlogMedia{},
		&models.Tag{},
	)
}

// Ping checks if database connection is alive
func Ping() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DBWithTimeout returns a new DB instance with timeout context
func DBWithTimeout(timeout time.Duration) *gorm.DB {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return DB.WithContext(ctx).Session(&gorm.Session{NewDB: true})
} 
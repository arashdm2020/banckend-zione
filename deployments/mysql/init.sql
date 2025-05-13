-- Create database if not exists
CREATE DATABASE IF NOT EXISTS zione_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use database
USE zione_db;

-- Create database structure
-- All tables are managed by GORM ORM, but we create some initial data here

-- Create initial roles
INSERT IGNORE INTO roles (id, name, created_at, updated_at) VALUES 
(1, 'admin', NOW(), NOW()),
(2, 'editor', NOW(), NOW()),
(3, 'user', NOW(), NOW());

-- Create initial admin user
-- Password: admin123 (hashed with bcrypt)
INSERT IGNORE INTO users (id, name, email, phone, password, role_id, created_at, updated_at) VALUES 
(1, 'Admin', 'admin@zionechain.cfd', '+989919980021', '$2a$10$6A.V5l/2bloUXj9FZ7VJVeI6lJkk4Ipz1C.cT7XBrDvL9zmz5nMy.', 1, NOW(), NOW());

-- Create initial blog categories
INSERT IGNORE INTO blog_categories (id, name, slug, created_at, updated_at) VALUES 
(1, 'Technology', 'technology', NOW(), NOW()),
(2, 'Design', 'design', NOW(), NOW()),
(3, 'Business', 'business', NOW(), NOW());

-- Create initial project categories
INSERT IGNORE INTO project_categories (id, name, slug, created_at, updated_at) VALUES 
(1, 'Web Development', 'web-development', NOW(), NOW()),
(2, 'Mobile Apps', 'mobile-apps', NOW(), NOW()),
(3, 'UI/UX Design', 'ui-ux-design', NOW(), NOW()); 
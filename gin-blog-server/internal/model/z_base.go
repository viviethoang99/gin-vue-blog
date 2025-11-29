package model

import (
	"time"

	"gorm.io/gorm"
)

// Migrate database tables. When there are no schema changes, it is recommended to comment this out.
// Supports creating tables and adding missing fields and indexes only.
// To protect data, it does not change existing field types or delete unused fields.
func MakeMigrate(db *gorm.DB) error {
	// Set up table associations
	db.SetupJoinTable(&Role{}, "Menus", &RoleMenu{})
	db.SetupJoinTable(&Role{}, "Resources", &RoleResource{})
	db.SetupJoinTable(&Role{}, "Users", &UserAuthRole{})
	db.SetupJoinTable(&UserAuth{}, "Roles", &UserAuthRole{})

	return db.AutoMigrate(
		&Article{},      // Article
		&Category{},     // Category
		&Tag{},          // Tag
		&Comment{},      // Comment
		&Message{},      // Message
		&FriendLink{},   // Friend Link
		&Page{},         // Page
		&Config{},       // Site Config
		&OperationLog{}, // Operation Log
		&UserInfo{},     // User Info

		&UserAuth{},     // User Auth
		&Role{},         // Role
		&Menu{},         // Menu
		&Resource{},     // Resource (API)
		&RoleMenu{},     // Role-Menu relation
		&RoleResource{}, // Role-Resource relation
		&UserAuthRole{}, // User-Role relation
	)
}

// Common model

type Model struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OptionVO struct {
	ID   int    `json:"value"`
	Name string `json:"label"`
}

// Gorm Scopes

// Pagination scope
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// Common CRUD

// Create data (single or batch)
func Create[T any](db *gorm.DB, data *T) (*T, error) {
	result := db.Create(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// Query single record
func Get[T any](db *gorm.DB, data *T, query string, args ...any) (*T, error) {
	result := db.Where(query, args...).First(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// Update single row: pass struct with primary key and struct with updated fields; zero values are not updated
func Update[T any](db *gorm.DB, data T, slt ...string) error {
	db = db.Model(&data)
	if len(slt) > 0 {
		db = db.Select(slt)
	}
	result := db.Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Batch update using map (map can update zero values); with conditions can perform single-row update
func UpdatesMap[T any](db *gorm.DB, data *T, maps map[string]any, query string, args ...any) error {
	result := db.Model(data).Where(query, args...).Updates(maps)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Batch update using struct fields (struct does not update zero values); with conditions can perform single-row update
func Updates[T any](db *gorm.DB, data T, query string, args ...any) error {
	result := db.Model(&data).Where(query, args...).Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// List data
func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// Batch delete; with conditions can delete a single record
func Delete[T any](db *gorm.DB, data T, query string, args ...any) error {
	result := db.Where(query, args...).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Count total
func Count[T any](db *gorm.DB, data *T, where ...any) (int, error) {
	var total int64
	db = db.Model(data)
	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}
	result := db.Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(total), nil
}

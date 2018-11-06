package main

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Model based in gorm.Model (Model)
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:",omitempty"`
}

// Author Struct (Model)
type Author struct {
	Model
	Firstname string // `json:"firstname"`
	Lastname  string // `json:"lastname"`
	Books     []Book `json:",omitempty"`
}

// Book Struct (Model)
type Book struct {
	Model
	Isbn     string  // `json:"isbn"`
	Title    string  // `json:"title"`
	Author   *Author `json:",omitempty"`
	AuthorID uint    `json:",omitempty"`
}

var db *gorm.DB

func initDB() error {
	DATABASE := os.Getenv("DATABASE")
	var err error
	if DATABASE == "" {
		db, err = gorm.Open("sqlite3", "/tmp/test.db")
	} else {
		// using gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
		db, err = gorm.Open("mysql", DATABASE)
	}
	if err != nil {
		return err
	}
	db.AutoMigrate(&Book{})
	db.AutoMigrate(&Author{})
	return nil
}

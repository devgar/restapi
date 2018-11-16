package main

import (
	"fmt"
	"os"
	"path"
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
	AuthorID uint    `json:",omitempty" sql:"type:integer REFERENCES authors(id) ON DELETE CASCADE ON UPDATE CASCADE"`
}

var db *gorm.DB

func scopedPath() string {
	return path.Join("/tmp", fmt.Sprintf("%d_%s", time.Now().Unix(), "test.db"))
}

func initDB() error {
	DATABASE := os.Getenv("DATABASE")
	var err error
	if DATABASE == "" {
		db, err = gorm.Open("sqlite3", scopedPath())
		db.Exec("PRAGMA foreign_keys = ON;")
	} else {
		// using gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
		db, err = gorm.Open("mysql", DATABASE)
	}
	if err != nil {
		return err
	}
	db.AutoMigrate(&Book{}, &Author{})
	return nil
}

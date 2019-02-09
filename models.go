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

// User Struct (Model)
type User struct {
	Model
	Username  string
	Email     string
	Firstname string
	Lastname  string
	Token     string
}

// BlogDomain Struct (Model)
type BlogDomain struct {
	Model
	Domain  string
	Owner   *User
	OwnerID uint `json:",omitempty" sql:"type:integer REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE"`
	Writers []User
}

// BlogPost Struct (Model)
type BlogPost struct {
	Model
	BlogDomain      *BlogDomain
	BlogDomainID    uint `json:",omitempty" sql:"type:integer REFERENCES blog_domains(id) ON DELETE SET NULL ON UPDATE CASCADE"`
	PublicationDate *time.Time
	Title           string
	Body            string
	BlogPostTags    []BlogPostTag
	Author          *User `json:",omitempty"`
	AuthorID        uint  `json:",omitempty" sql:"type:integer REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE"`
}

// BlogPostTag Struct (Model)
type BlogPostTag struct {
	Model
	BlogDomain   *BlogDomain
	BlogDomainID uint `json:",omitempty" sql:"type:integer REFERENCES blog_domains(id) ON DELETE SET NULL ON UPDATE CASCADE"`
	Tag          string
	BlogPosts    []BlogPost
}

// BlogPostCategory (Model)
type BlogPostCategory struct {
	Model
	BlogDomain       *BlogDomain
	BlogDomainID     uint `json:",omitempty" sql:"type:integer REFERENCES blog_domains(id) ON DELETE SET NULL ON UPDATE CASCADE"`
	Category         string
	BlogPosts        []BlogPost
	ParentCategory   *BlogPostCategory `json:",omitempty"`
	ParentCategoryID uint              `json:",omitempty" sql:"type:integer REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE"`
}

var db *gorm.DB

func sqliteDBPath() string {
	TMPDIR := os.Getenv("TMPDIR")
	if TMPDIR == "" {
		TMPDIR = "/tmp"
	}
	return path.Join(TMPDIR, fmt.Sprintf("%d_%s", time.Now().Unix(), "test.db"))
}

func initDB() error {
	DATABASE := os.Getenv("DATABASE")
	var err error
	if DATABASE == "" {
		db, err = gorm.Open("sqlite3", sqliteDBPath())
		db.Exec("PRAGMA foreign_keys = ON;")
	} else {
		// using gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
		db, err = gorm.Open("mysql", DATABASE)
	}
	if err != nil {
		return err
	}
	db.AutoMigrate(&Book{}, &Author{})
	db.AutoMigrate(&BlogDomain{}, &User{}, &BlogPost{}, &BlogPostCategory{}, &BlogPostTag{})
	return nil
}

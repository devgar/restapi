package database

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/devgar/restapi/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // configures mysql driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func sqliteDBPath() string {
	TMPDIR := os.Getenv("TMPDIR")
	if TMPDIR == "" {
		TMPDIR = "/tmp"
	}
	return path.Join(TMPDIR, fmt.Sprintf("%d_%s", time.Now().Unix(), "test.db"))
}

// Initialize initializes the database
func Initialize() (*gorm.DB, error) {
	DATABASE := os.Getenv("DATABASE")
	var err error
	var db *gorm.DB
	if DATABASE == "" {
		db, err = gorm.Open("sqlite3", sqliteDBPath())
		db.Exec("PRAGMA foreign_keys = ON;")
	} else {
		// using gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
		db, err = gorm.Open("mysql", DATABASE)
	}
	if err != nil {
		return nil, err
	}
	db.LogMode(true) // logs SQL
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	models.Migrate(db)
	return db, err
}

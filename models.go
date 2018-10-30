package main

import "time"

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

// DetailBook (Model) // * Ignores AuthorID
type DetailBook struct {
	Model
	Isbn   string
	Title  string
	Author *Author `json:",omitempty"`
}

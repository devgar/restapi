package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeBlogPostsGet(c *gin.Context) {
	var blogPosts []BlogPost
	db.Find(&blogPosts)
	c.JSON(200, blogPosts)
}

func routeBlogPostsGetOne(c *gin.Context) {
	id := c.Param("id")
	var blogPost BlogPost
	db.Preload("Author").First(&blogPost, id)
	if blogPost.ID == 0 {
		c.JSON(404, "BlogPost not found")
	} else {
		blogPost.AuthorID = 0
		c.JSON(200, blogPost)
	}
}

func routeBlogPostsPost(c *gin.Context) {
	var blogPost BlogPost
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&blogPost).Error; err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blogPost)
}

func routeBlogPostsPut(c *gin.Context) {
	id := c.Param("id")
	var changes, blogPost BlogPost
	if err := c.ShouldBindJSON(&changes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.First(&blogPost, id)
	if blogPost.ID == 0 {
		c.JSON(404, "BlogPost not found")
		return
	}
	if err := db.Model(&blogPost).Updates(&changes).Error; err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, blogPost)
}

func routeBlogPostsDelete(c *gin.Context) {
	id := c.Param("id")
	var blogPost BlogPost
	db.First(&blogPost, id)
	if blogPost.ID == 0 {
		c.JSON(404, "Book not found")
	} else {
		db.Delete(&blogPost)
		c.JSON(200, blogPost)
	}
}

func routesBlogPosts(r *gin.RouterGroup) {
	r.GET("", routeBlogPostsGet)
	r.GET("/:id", routeBlogPostsGetOne)
	r.POST("", routeBlogPostsPost)
	r.PUT("/:id", routeBlogPostsPut)
	r.DELETE("/:id", routeBlogPostsDelete)
}

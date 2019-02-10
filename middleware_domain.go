package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var domainKeys map[string]int

func middlewareDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("x-dk")
		if domain, ok := domainKeys[key]; !ok {
			panic("Error 'x-dk' not found")
		} else {
			c.Set("domain", domain)
		}
		t := time.Now()
		c.Next()
		log.Print(time.Since(t))
		log.Println(c.Writer.Status())
	}
}

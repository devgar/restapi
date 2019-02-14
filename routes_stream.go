package main

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func routesStream(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("X-Accel-Buffering", "no")

		rw := c.Writer

		flusher, ok := rw.(http.Flusher)
		if !ok {
			c.String(500, "Streaming not supported")
			return
		}

		io.WriteString(rw, ": ping\n\n")
		flusher.Flush()

		eventc := make(chan int, 10)

		rc := c.Request.Context()
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		go func() { // TODO: goroutine for developing
			ticker := time.NewTicker(3 * time.Second)
			ender := time.After(time.Minute)
			defer close(eventc)
			// defer resolve()
			for {
				select {
				case <-ender:
					return
				case <-ticker.C:
					eventc <- rand.Intn(100)
				case <-ctx.Done():
					return
				}
			}
		}()

		id := 0
		pinger := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-rc.Done():
				return
			case <-c.Done():
				return
			case <-ctx.Done():
				return
			case <-pinger.C:
				if c.IsAborted() {
					c.Abort()
					return
				}
				io.WriteString(rw, ": ping\n\n")
				flusher.Flush()
			case data, ok := <-eventc:
				if ok {
					io.WriteString(rw, "id: "+strconv.Itoa(id))
					io.WriteString(rw, " ")
					io.WriteString(rw, "data: "+strconv.Itoa(data))
					io.WriteString(rw, "\n\n")
					flusher.Flush()
				} else {
					return
				}
				id++
			}

		}
	})
}

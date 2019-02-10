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
	r.GET("/stream", func(c *gin.Context) {
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
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer close(eventc)
			// defer resolve()
			for {
				select {
				case <-time.After(5 * time.Minute):
					return
				case <-ticker.C:
					eventc <- rand.Intn(100)
				}
			}

		}()

		id := 0

		for {
			select {
			case <-time.After(time.Hour):
				return
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 30):
				io.WriteString(rw, ": ping\n\n")
				flusher.Flush()
			case buf, ok := <-eventc:
				if ok {
					io.WriteString(rw, "id: "+strconv.Itoa(id))
					io.WriteString(rw, "\n")
					io.WriteString(rw, "data: "+strconv.Itoa(buf))
					io.WriteString(rw, "\n\n")
					flusher.Flush()
				}
				id++
			}

		}
	})
}

package main

import (
	"context"
	"fmt"
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
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		go func() { // goroutine
			ticker := time.NewTicker(5 * time.Second)
			ender := time.After(time.Minute)
			defer close(eventc)
			// defer resolve()
			for {
				select {
				case <-ender:
					fmt.Println("|| -- Breaking -- ||")
					return
				case <-ticker.C:
					eventc <- rand.Intn(100)
				}
			}
		}()

		id := 0
		pinger := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-c.Done():
				fmt.Println("|| -- (C) Context Done --||")
				return
			case <-ctx.Done():
				fmt.Println("|| -- Context Done --||")
				return
			case <-pinger.C:
				if c.IsAborted() {
					fmt.Println("|| -- Aborted -- ||")
					c.Abort()
					return
				}
				io.WriteString(rw, ": ping\n\n")
				flusher.Flush()
			case data, ok := <-eventc:
				if err := ctx.Err(); err != nil {
					fmt.Println("|| -- ctx:Error -- ||")
					return
				}
				if ok {
					io.WriteString(rw, "id: "+strconv.Itoa(id))
					io.WriteString(rw, "\n")
					io.WriteString(rw, "data: "+strconv.Itoa(data))
					io.WriteString(rw, "\n\n")
					flusher.Flush()
				} else {
					fmt.Println("|| -- eventC CLOSED -- ||")
					// c.Abort()
					return
				}
				id++
			}

		}
	})
}

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

var registeredBookChannels []chan Book

func unregistBookChannel(c chan Book) {
	for i, v := range registeredBookChannels {
		if v == c {
			registeredBookChannels[i] = registeredBookChannels[len(registeredBookChannels)-1]
			registeredBookChannels = registeredBookChannels[:len(registeredBookChannels)-1]
			break
		}
	}
}

var domainChannels []chan StreamData

func unregistDomainChannel(c chan StreamData) {
	for i, v := range domainChannels {
		if v == c {
			domainChannels[i] = domainChannels[len(domainChannels)-1]
			domainChannels = domainChannels[:len(domainChannels)-1]
			break
		}
	}
}

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

	r.GET("/books", func(c *gin.Context) {
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

		var ping = func() {
			io.WriteString(rw, ": ping ")
			io.WriteString(rw, strconv.FormatInt(time.Now().Unix(), 10))
			io.WriteString(rw, "\n\n")
			flusher.Flush()
		}

		ping()

		eventc := make(chan Book)
		registeredBookChannels = append(registeredBookChannels, eventc)

		rc := c.Request.Context()
		ctx, cancel := context.WithCancel(context.Background())

		defer func() {
			unregistBookChannel(eventc)
			close(eventc)
			cancel()
		}()
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
				ping()
			case book, ok := <-eventc:
				if ok {
					io.WriteString(rw, "POST book: "+strconv.Itoa(int(book.ID)))
					io.WriteString(rw, "\n\n")
					flusher.Flush()
				} else {
					return
				}
			}

		}
	})
}

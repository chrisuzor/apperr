// Package apperrgin classifies an error and hands the status, code,
// client-safe message, and original error to a renderer.
// The application writes the response body.
package apperrgin

import (
	"github.com/chrisuzor/apperr/apperrhttp"
	"github.com/gin-gonic/gin"
)

// Renderer writes the HTTP response. The application owns the JSON shape.
// message is safe for the client. err is the original error, for logs and for
// reading a typed cause such as a validation error.
type Renderer func(c *gin.Context, status int, code, message string, err error)

// Writer calls a renderer stored at construction.
type Writer struct {
	render Renderer
}

// New returns a writer that uses render for every error.
// It panics if render is nil.
func New(render Renderer) *Writer {
	if render == nil {
		panic("apperrgin: nil renderer")
	}
	return &Writer{render: render}
}

// Write classifies err, records it on the Gin context, stops the chain, and calls render.
// A nil context or a nil error does not write a response.
// A nil renderer panics when err is non-nil.
func (w *Writer) Write(c *gin.Context, err error) {
	if c == nil || err == nil {
		return
	}
	if w == nil || w.render == nil {
		panic("apperrgin: nil renderer")
	}
	write(c, err, w.render)
}

// Write classifies err, records it on the Gin context, stops the chain, and calls render.
// A nil context or a nil error does not write a response.
// A nil renderer panics when err is non-nil.
func Write(c *gin.Context, err error, render Renderer) {
	write(c, err, render)
}

func write(c *gin.Context, err error, render Renderer) {
	if c == nil || err == nil {
		return
	}
	if render == nil {
		panic("apperrgin: nil renderer")
	}
	_ = c.Error(err)
	c.Abort()
	render(c, apperrhttp.Status(err), apperrhttp.Code(err), apperrhttp.PublicMessage(err), err)
}

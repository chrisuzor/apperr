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
func New(render Renderer) *Writer {
	return &Writer{render: render}
}

// Write classifies err, records it on the Gin context, stops the chain, and calls render.
// A nil context, a nil error, or a nil renderer does not write a response.
func (w *Writer) Write(c *gin.Context, err error) {
	if w == nil {
		return
	}
	write(c, err, w.render)
}

// Write classifies err, records it on the Gin context, stops the chain, and calls render.
// A nil context, a nil error, or a nil renderer does not write a response.
func Write(c *gin.Context, err error, render Renderer) {
	write(c, err, render)
}

func write(c *gin.Context, err error, render Renderer) {
	if c == nil || err == nil || render == nil {
		return
	}
	_ = c.Error(err)
	c.Abort()
	render(c, apperrhttp.Status(err), apperrhttp.Code(err), apperrhttp.PublicMessage(err), err)
}

// Package apperrgin classifies an error and hands the status, code, and
// client-safe message to a renderer. The application writes the response body.
package apperrgin

import (
	"github.com/chrisuzor/apperr/httpx"
	"github.com/gin-gonic/gin"
)

// Renderer writes the HTTP response. The application owns the JSON shape.
type Renderer func(c *gin.Context, status int, code, message string)

// Write classifies err and calls render.
// A nil context, a nil error, or a nil renderer does not write a response.
func Write(c *gin.Context, err error, render Renderer) {
	if c == nil || err == nil || render == nil {
		return
	}
	render(c, httpx.Status(err), httpx.Code(err), httpx.PublicMessage(err))
}

// Package httpx maps an apperr code to an HTTP status and a client-safe message.
// It does not write a response.
package httpx

import (
	"net/http"

	"github.com/chrisuzor/apperr"
)

const internalMessage = "internal error"

// Status returns the HTTP status for err.
// [apperr.CodeInvalid] is 400, [apperr.CodeNotFound] is 404,
// [apperr.CodeForbidden] is 403, and [apperr.CodeConflict] is 409.
// Every other code, including a plain error, is 500.
func Status(err error) int {
	switch apperr.CodeOf(err) {
	case apperr.CodeInvalid:
		return http.StatusBadRequest
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeForbidden:
		return http.StatusForbidden
	case apperr.CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// PublicMessage returns text that is safe to send to a client.
// Internal and unclassified errors return a fixed message. The wrapped cause
// stays on err for logs.
func PublicMessage(err error) string {
	if apperr.CodeOf(err) == apperr.CodeInternal {
		return internalMessage
	}
	if msg := apperr.Message(err); msg != "" {
		return msg
	}
	return http.StatusText(Status(err))
}

// Code returns the apperr code as a string.
func Code(err error) string {
	return string(apperr.CodeOf(err))
}

// Package apperrhttp maps an apperr code to an HTTP status and a client-safe message.
// It does not write a response.
package apperrhttp

import (
	"net/http"

	"github.com/chrisuzor/apperr"
)

const internalMessage = "internal error"

// Status returns the HTTP status for err.
// [apperr.CodeInvalid] is 400, [apperr.CodeUnauthenticated] is 401,
// [apperr.CodeForbidden] is 403, [apperr.CodeNotFound] is 404,
// [apperr.CodeConflict] is 409, [apperr.CodeRateLimited] is 429,
// and [apperr.CodeUnavailable] is 503.
// Every other code, including a plain error, is 500.
func Status(err error) int {
	switch apperr.CodeOf(err) {
	case apperr.CodeInvalid:
		return http.StatusBadRequest
	case apperr.CodeUnauthenticated:
		return http.StatusUnauthorized
	case apperr.CodeForbidden:
		return http.StatusForbidden
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeConflict:
		return http.StatusConflict
	case apperr.CodeRateLimited:
		return http.StatusTooManyRequests
	case apperr.CodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// PublicMessage returns text that is safe to send to a client.
// A 500 returns a fixed message, including when the code is not one this
// package knows. The wrapped cause stays on err for logs.
func PublicMessage(err error) string {
	if Status(err) == http.StatusInternalServerError {
		return internalMessage
	}
	if msg := apperr.Message(err); msg != "" {
		return msg
	}
	return http.StatusText(Status(err))
}

// Code returns the code a client may see.
// A code this package does not map, and [apperr.CodeInternal], are reported as internal.
func Code(err error) string {
	switch apperr.CodeOf(err) {
	case apperr.CodeInvalid, apperr.CodeUnauthenticated, apperr.CodeForbidden,
		apperr.CodeNotFound, apperr.CodeConflict, apperr.CodeRateLimited,
		apperr.CodeUnavailable:
		return string(apperr.CodeOf(err))
	default:
		return string(apperr.CodeInternal)
	}
}

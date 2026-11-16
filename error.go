package apperr

import "errors"

// Code is the category an edge uses to choose a response.
type Code string

const (
	CodeInvalid   Code = "invalid"
	CodeNotFound  Code = "not_found"
	CodeForbidden Code = "forbidden"
	CodeConflict  Code = "conflict"
	CodeInternal  Code = "internal"
)

// Error is an application error with a category and a client-safe message.
// The wrapped cause is for logs. Edges that speak to clients leave it off the response.
type Error struct {
	code Code
	msg  string
	err  error
}

// New returns an error with the given code and client-safe message.
// An empty code is stored as [CodeInternal].
func New(code Code, msg string) error {
	if code == "" {
		code = CodeInternal
	}
	return &Error{code: code, msg: msg}
}

// Wrap returns an error that classifies as code and unwraps to err.
// A nil err returns nil. An empty code is stored as [CodeInternal].
func Wrap(code Code, err error, msg string) error {
	if err == nil {
		return nil
	}
	if code == "" {
		code = CodeInternal
	}
	return &Error{code: code, msg: msg, err: err}
}

// Replace returns replacement when errors.Is(err, match).
// A nil err, a nil match, or a nil replacement returns err unchanged.
func Replace(err, match, replacement error) error {
	if err == nil || match == nil || replacement == nil {
		return err
	}
	if errors.Is(err, match) {
		return replacement
	}
	return err
}

// CodeOf returns the code of the outermost [*Error] in err.
// It returns [CodeInternal] when err is nil or contains no [*Error].
func CodeOf(err error) Code {
	e, ok := As(err)
	if !ok || e.code == "" {
		return CodeInternal
	}
	return e.code
}

// Message returns the client-safe message of the outermost [*Error].
// It returns an empty string when err contains no [*Error].
// The message is whatever the caller stored, including messages on internal errors.
// Use an edge helper when the client must not see internal text.
func Message(err error) string {
	e, ok := As(err)
	if !ok {
		return ""
	}
	return e.msg
}

// As finds the outermost [*Error] in err.
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// Error returns the message and, when present, the wrapped cause.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.err == nil {
		return e.msg
	}
	if e.msg == "" {
		return e.err.Error()
	}
	return e.msg + ": " + e.err.Error()
}

// Unwrap returns the wrapped cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// Code returns the category. A nil receiver returns [CodeInternal].
func (e *Error) Code() Code {
	if e == nil || e.code == "" {
		return CodeInternal
	}
	return e.code
}

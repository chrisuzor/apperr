// Package apperr classifies errors at an application boundary.
//
// A service returns an [*Error] with a [Code]. An edge such as HTTP maps that
// code to a status. This package does not choose a response body.
package apperr

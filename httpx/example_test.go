package httpx_test

import (
	"errors"
	"fmt"

	"github.com/chrisuzor/apperr"
	"github.com/chrisuzor/apperr/httpx"
)

func ExamplePublicMessage() {
	err := apperr.Wrap(apperr.CodeInternal, errors.New("dial tcp: connection refused"), "load question")
	fmt.Println(httpx.Status(err))
	fmt.Println(httpx.PublicMessage(err))
	fmt.Println(err)
	// Output:
	// 500
	// internal error
	// load question: dial tcp: connection refused
}

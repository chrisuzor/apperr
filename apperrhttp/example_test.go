package apperrhttp_test

import (
	"errors"
	"fmt"

	"github.com/chrisuzor/apperr"
	"github.com/chrisuzor/apperr/apperrhttp"
)

func ExamplePublicMessage() {
	err := apperr.Wrap(apperr.CodeInternal, errors.New("dial tcp: connection refused"), "load question")
	fmt.Println(apperrhttp.Status(err))
	fmt.Println(apperrhttp.PublicMessage(err))
	fmt.Println(err)
	// Output:
	// 500
	// internal error
	// load question: dial tcp: connection refused
}

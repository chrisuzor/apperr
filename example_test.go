package apperr_test

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chrisuzor/apperr"
)

func ExampleReplace() {
	err := apperr.Replace(sql.ErrNoRows, sql.ErrNoRows, apperr.New(apperr.CodeNotFound, "question not found"))
	fmt.Println(apperr.CodeOf(err))
	fmt.Println(apperr.Message(err))
	// Output:
	// not_found
	// question not found
}

func ExampleWrap() {
	err := apperr.Wrap(apperr.CodeInternal, errors.New("dial tcp: connection refused"), "load question")
	fmt.Println(apperr.CodeOf(err))
	fmt.Println(err)
	// Output:
	// internal
	// load question: dial tcp: connection refused
}

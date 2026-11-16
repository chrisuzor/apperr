package httpx

import (
	"errors"
	"net/http"
	"testing"

	"github.com/chrisuzor/apperr"
)

func TestStatus(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want int
	}{
		{"invalid", apperr.New(apperr.CodeInvalid, "bad id"), http.StatusBadRequest},
		{"not found", apperr.New(apperr.CodeNotFound, "question not found"), http.StatusNotFound},
		{"forbidden", apperr.New(apperr.CodeForbidden, "exam started"), http.StatusForbidden},
		{"conflict", apperr.New(apperr.CodeConflict, "already submitted"), http.StatusConflict},
		{"internal", apperr.Wrap(apperr.CodeInternal, errors.New("dial tcp"), "load question"), http.StatusInternalServerError},
		{"plain", errors.New("boom"), http.StatusInternalServerError},
		{"nil", nil, http.StatusInternalServerError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := Status(tc.err); got != tc.want {
				t.Fatalf("Status = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestPublicMessageHidesInternalCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("pq: duplicate key value violates unique constraint")
	err := apperr.Wrap(apperr.CodeInternal, cause, "save question")

	if got := PublicMessage(err); got != "internal error" {
		t.Fatalf("PublicMessage = %q", got)
	}
	if got := err.Error(); got != "save question: "+cause.Error() {
		t.Fatalf("Error = %q", got)
	}

	visible := apperr.New(apperr.CodeNotFound, "question not found")
	if got := PublicMessage(visible); got != "question not found" {
		t.Fatalf("PublicMessage = %q", got)
	}
	if got := Code(visible); got != string(apperr.CodeNotFound) {
		t.Fatalf("Code = %q", got)
	}
}

func TestPublicMessageEmptyUsesStatusText(t *testing.T) {
	t.Parallel()

	err := apperr.New(apperr.CodeConflict, "")
	if got := PublicMessage(err); got != http.StatusText(http.StatusConflict) {
		t.Fatalf("PublicMessage = %q", got)
	}
}

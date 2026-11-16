package apperrhttp

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
		{"unauthenticated", apperr.New(apperr.CodeUnauthenticated, "sign in"), http.StatusUnauthorized},
		{"forbidden", apperr.New(apperr.CodeForbidden, "exam started"), http.StatusForbidden},
		{"not found", apperr.New(apperr.CodeNotFound, "question not found"), http.StatusNotFound},
		{"conflict", apperr.New(apperr.CodeConflict, "already submitted"), http.StatusConflict},
		{"rate limited", apperr.New(apperr.CodeRateLimited, "slow down"), http.StatusTooManyRequests},
		{"unavailable", apperr.New(apperr.CodeUnavailable, "try later"), http.StatusServiceUnavailable},
		{"internal", apperr.Wrap(apperr.CodeInternal, errors.New("dial tcp"), "load question"), http.StatusInternalServerError},
		{"plain", errors.New("boom"), http.StatusInternalServerError},
		{"nil", nil, http.StatusInternalServerError},
		{"unknown", apperr.New("teapot", "secret detail"), http.StatusInternalServerError},
	}
	for _, tc := range cases {
		tc := tc
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

func TestUnknownCodeDoesNotLeak(t *testing.T) {
	t.Parallel()

	err := apperr.New("teapot", "secret detail")
	if got := PublicMessage(err); got != "internal error" {
		t.Fatalf("PublicMessage = %q", got)
	}
	if got := Code(err); got != string(apperr.CodeInternal) {
		t.Fatalf("Code = %q", got)
	}
	if apperr.CodeOf(err) != "teapot" {
		t.Fatalf("CodeOf = %q", apperr.CodeOf(err))
	}
}

func TestPublicMessageEmptyUsesStatusText(t *testing.T) {
	t.Parallel()

	err := apperr.New(apperr.CodeConflict, "")
	if got := PublicMessage(err); got != http.StatusText(http.StatusConflict) {
		t.Fatalf("PublicMessage = %q", got)
	}
}

func TestWrapInternalHidesEarlierCode(t *testing.T) {
	t.Parallel()

	err := apperr.Wrap(apperr.CodeInternal, apperr.New(apperr.CodeNotFound, "question not found"), "load question")
	if Status(err) != http.StatusInternalServerError {
		t.Fatalf("Status = %d", Status(err))
	}
	if PublicMessage(err) != "internal error" {
		t.Fatalf("PublicMessage = %q", PublicMessage(err))
	}
}

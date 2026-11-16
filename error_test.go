package apperr

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func TestNewAndWrap(t *testing.T) {
	t.Parallel()

	e := New(CodeNotFound, "question not found")
	if CodeOf(e) != CodeNotFound {
		t.Fatalf("CodeOf = %q", CodeOf(e))
	}
	if Message(e) != "question not found" {
		t.Fatalf("Message = %q", Message(e))
	}
	if e.Error() != "question not found" {
		t.Fatalf("Error = %q", e.Error())
	}

	cause := errors.New("dial tcp")
	wrapped := Wrap(CodeInternal, cause, "load question")
	if !errors.Is(wrapped, cause) {
		t.Fatal("Wrap did not keep the cause")
	}
	if wrapped.Error() != "load question: dial tcp" {
		t.Fatalf("Error = %q", wrapped.Error())
	}
	if Wrap(CodeNotFound, nil, "ignored") != nil {
		t.Fatal("Wrap(nil) returned an error")
	}
	if Annotate(nil, "load question") != nil {
		t.Fatal("Annotate(nil) returned an error")
	}
	if CodeOf(New("", "x")) != CodeInternal {
		t.Fatal("empty code was not stored as internal")
	}
}

func TestCodeOfThroughWrapper(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("handler: %w", New(CodeForbidden, "exam started"))
	if CodeOf(err) != CodeForbidden {
		t.Fatalf("CodeOf = %q", CodeOf(err))
	}
	if Message(err) != "exam started" {
		t.Fatalf("Message = %q", Message(err))
	}
	if CodeOf(errors.New("plain")) != CodeInternal {
		t.Fatal("plain error was not internal")
	}
	if CodeOf(nil) != CodeInternal {
		t.Fatal("nil was not internal")
	}
	if Message(nil) != "" {
		t.Fatal("nil message was not empty")
	}
}

func TestReplace(t *testing.T) {
	t.Parallel()

	notFound := New(CodeNotFound, "question not found")

	got := Replace(sql.ErrNoRows, sql.ErrNoRows, notFound)
	if got != notFound {
		t.Fatalf("Replace matched sentinel, got %v", got)
	}

	other := errors.New("timeout")
	if Replace(other, sql.ErrNoRows, notFound) != other {
		t.Fatal("Replace changed an unmatched error")
	}
	if Replace(nil, sql.ErrNoRows, notFound) != nil {
		t.Fatal("Replace changed a nil error")
	}
	if Replace(sql.ErrNoRows, nil, notFound) != sql.ErrNoRows {
		t.Fatal("Replace matched a nil sentinel")
	}
	if Replace(sql.ErrNoRows, sql.ErrNoRows, nil) != sql.ErrNoRows {
		t.Fatal("Replace used a nil replacement")
	}

	wrapped := fmt.Errorf("query: %w", sql.ErrNoRows)
	if Replace(wrapped, sql.ErrNoRows, notFound) != notFound {
		t.Fatal("Replace did not use errors.Is")
	}
}

func TestAnnotateKeepsCode(t *testing.T) {
	t.Parallel()

	inner := New(CodeNotFound, "question not found")
	got := Annotate(inner, "load question")
	if CodeOf(got) != CodeNotFound {
		t.Fatalf("CodeOf = %q", CodeOf(got))
	}
	if Message(got) != "question not found" {
		t.Fatalf("Message = %q", Message(got))
	}
	if got.Error() != "load question: question not found" {
		t.Fatalf("Error = %q", got.Error())
	}
	if Annotate(inner, "") != inner {
		t.Fatal("Annotate with an empty message changed the error")
	}
}

func TestReplaceOptInKeepsCause(t *testing.T) {
	t.Parallel()

	replacement := Wrap(CodeNotFound, sql.ErrNoRows, "question not found")
	got := Replace(sql.ErrNoRows, sql.ErrNoRows, replacement)
	if CodeOf(got) != CodeNotFound {
		t.Fatalf("CodeOf = %q", CodeOf(got))
	}
	if !errors.Is(got, sql.ErrNoRows) {
		t.Fatal("opt-in replacement dropped sql.ErrNoRows")
	}

	plain := Replace(sql.ErrNoRows, sql.ErrNoRows, New(CodeNotFound, "question not found"))
	if errors.Is(plain, sql.ErrNoRows) {
		t.Fatal("default replacement kept sql.ErrNoRows")
	}
}

func TestOutermostCodeWins(t *testing.T) {
	t.Parallel()

	inner := New(CodeNotFound, "missing")
	outer := Wrap(CodeConflict, inner, "duplicate")
	if CodeOf(outer) != CodeConflict {
		t.Fatalf("CodeOf = %q", CodeOf(outer))
	}
	if Message(outer) != "duplicate" {
		t.Fatalf("Message = %q", Message(outer))
	}
}

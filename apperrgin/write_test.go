package apperrgin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chrisuzor/apperr"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestWrite(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := apperr.New(apperr.CodeForbidden, "exam started")

	called := false
	New(func(_ *gin.Context, status int, code, message string, got error) {
		called = true
		if status != http.StatusForbidden || code != "forbidden" || message != "exam started" || got != err {
			t.Fatalf("status=%d code=%q message=%q err=%v", status, code, message, got)
		}
	}).Write(c, err)
	if !called {
		t.Fatal("renderer was not called")
	}
	if !c.IsAborted() {
		t.Fatal("context was not aborted")
	}
	if len(c.Errors) != 1 || c.Errors[0].Err != err {
		t.Fatalf("context errors = %v", c.Errors)
	}
}

func TestWriteHidesInternalMessage(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	err := apperr.Wrap(apperr.CodeInternal, errors.New("relation exam does not exist"), "load exam")

	called := false
	Write(c, err, func(_ *gin.Context, status int, code, message string, got error) {
		called = true
		if status != http.StatusInternalServerError || code != "internal" || message != "internal error" || got != err {
			t.Fatalf("status=%d code=%q message=%q err=%v", status, code, message, got)
		}
	})
	if !called {
		t.Fatal("renderer was not called")
	}
}

func TestWriteSkipsNil(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Write(c, nil, func(*gin.Context, int, string, string, error) {
		t.Fatal("rendered a nil error")
	})
	Write(nil, apperr.New(apperr.CodeNotFound, "missing"), func(*gin.Context, int, string, string, error) {
		t.Fatal("rendered with a nil context")
	})
	Write(c, nil, nil)

	if w.Body.Len() != 0 {
		t.Fatalf("wrote a body: %s", w.Body.String())
	}
	if c.IsAborted() {
		t.Fatal("aborted a nil error")
	}
	if len(c.Errors) != 0 {
		t.Fatalf("recorded errors for a nil error: %v", c.Errors)
	}
}

func TestNilRendererPanics(t *testing.T) {
	t.Parallel()

	err := apperr.New(apperr.CodeNotFound, "missing")

	t.Run("constructor", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("New(nil) did not panic")
			}
		}()
		New(nil)
	})

	t.Run("write function", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		defer func() {
			if recover() == nil {
				t.Fatal("Write with a nil renderer did not panic")
			}
		}()
		Write(c, err, nil)
	})
}

func TestNilWriterPanics(t *testing.T) {
	t.Parallel()

	err := apperr.New(apperr.CodeNotFound, "missing")

	t.Run("nil pointer", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		defer func() {
			if recover() == nil {
				t.Fatal("nil Writer did not panic")
			}
		}()
		(*Writer)(nil).Write(c, err)
	})

	t.Run("zero value", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		defer func() {
			if recover() == nil {
				t.Fatal("zero Writer did not panic")
			}
		}()
		(&Writer{}).Write(c, err)
	})
}

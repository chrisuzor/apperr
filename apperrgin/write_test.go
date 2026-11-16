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

	called := false
	Write(c, apperr.New(apperr.CodeForbidden, "exam started"), func(_ *gin.Context, status int, code, message string) {
		called = true
		if status != http.StatusForbidden || code != "forbidden" || message != "exam started" {
			t.Fatalf("status=%d code=%q message=%q", status, code, message)
		}
	})
	if !called {
		t.Fatal("renderer was not called")
	}
}

func TestWriteHidesInternalMessage(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	err := apperr.Wrap(apperr.CodeInternal, errors.New("relation exam does not exist"), "load exam")
	Write(c, err, func(_ *gin.Context, status int, code, message string) {
		if status != http.StatusInternalServerError || code != "internal" || message != "internal error" {
			t.Fatalf("status=%d code=%q message=%q", status, code, message)
		}
	})
}

func TestWriteSkipsNil(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Write(c, nil, func(*gin.Context, int, string, string) {
		t.Fatal("rendered a nil error")
	})
	Write(nil, apperr.New(apperr.CodeNotFound, "missing"), func(*gin.Context, int, string, string) {
		t.Fatal("rendered with a nil context")
	})
	Write(c, apperr.New(apperr.CodeNotFound, "missing"), nil)

	if w.Body.Len() != 0 {
		t.Fatalf("wrote a body: %s", w.Body.String())
	}
}

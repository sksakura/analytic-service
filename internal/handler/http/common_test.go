package handlerhttp

import (
	"analytic-service/internal/db"
	"analytic-service/internal/logger"
	"analytic-service/internal/profiler"
	"context"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestSayHello(t *testing.T) {
	t.Run("hello_without_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		mockedLogger := logger.Mock(zaptest.NewLogger(t))
		New(
			db.NewReport(
				context.Background(),
				"MEM",
				"mem",
				mockedLogger),
			&profiler.Profiler{},
			mockedLogger,
			&Auth{
				Tokensservice: nil,
			}).SayHello(w, r)
		body := w.Body.String()
		if body != "welcome" {
			t.Error(body, "welcome")
		}
	})
}

func TestGetPprofStatus(t *testing.T) {
	t.Run("GetPprofStatus", func(t *testing.T) {
		mockedLogger := logger.Mock(zaptest.NewLogger(t))
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		New(
			db.NewReport(
				context.Background(),
				"MEM",
				"mem",
				mockedLogger),
			profiler.NewProfiler(false),
			mockedLogger,
			&Auth{
				Tokensservice: nil,
			}).GetPprof(w, r)

		body := w.Body.String()
		exp := "pprof_on state:false"
		if body != exp {
			t.Error(body, exp)
		}
	})
}

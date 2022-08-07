package handlerhttp

import (
	"analytic-service/internal/db"
	"analytic-service/internal/logger"
	"analytic-service/internal/profiler"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/g6834/team21/authproto.git/pkg/tokensservice"
	"google.golang.org/grpc"
)

type MockedTokensServiceClient struct {
}

func (MockedTokensServiceClient) Validate(ctx context.Context, in *tokensservice.Tokens, opts ...grpc.CallOption) (*tokensservice.ValidationResult, error) {
	return &tokensservice.ValidationResult{UserID: "mockedUser"}, nil
}

var mockedLogger = logger.Mock(nil)
var rpt = db.NewReport(context.Background(), "MEM", "memory", mockedLogger)
var mockedReportHandler = New(
	rpt,
	&profiler.Profiler{},
	mockedLogger,
	&Auth{
		Tokensservice: MockedTokensServiceClient{},
	})

func Test_ReportTaskList(t *testing.T) {
	t.Run("tasklist_no_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		mockedReportHandler.ReportTaskList(w, r)
		status := w.Result().StatusCode
		if status != http.StatusUnauthorized {
			t.Error(status, http.StatusUnauthorized)
		}
	})

	t.Run("tasklist_with_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.AddCookie(&http.Cookie{
			Name:     "access_token",
			Value:    "12345",
			HttpOnly: true,
		})
		r.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "-12345",
			HttpOnly: true,
		})
		nextHandler := http.HandlerFunc(mockedReportHandler.ReportTaskList)
		handlerToTest := mockedReportHandler.AuthMiddleware(nextHandler)

		handlerToTest.ServeHTTP(w, r)

		status := w.Result().StatusCode
		if status != http.StatusOK {
			t.Error(status, http.StatusOK)
		}
		expected := `{"type":"livetime","content":"[]"}`
		body := w.Body.String()
		if body != expected {
			t.Error(body, expected)
		}
	})
}

func TestReportApproved(t *testing.T) {
	t.Run("approved_no_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		mockedReportHandler.ReportApproved(w, r)
		status := w.Result().StatusCode
		if status != http.StatusUnauthorized {
			t.Error(status, http.StatusUnauthorized)
		}
	})

	t.Run("approved_with_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.AddCookie(&http.Cookie{
			Name:     "access_token",
			Value:    "12345",
			HttpOnly: true,
		})
		r.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "-12345",
			HttpOnly: true,
		})

		nextHandler := http.HandlerFunc(mockedReportHandler.ReportApproved)
		handlerToTest := mockedReportHandler.AuthMiddleware(nextHandler)

		handlerToTest.ServeHTTP(w, r)

		status := w.Result().StatusCode
		if status != http.StatusOK {
			t.Error(status, http.StatusOK)
		}
		expected := `{"type":"approved","content":"111"}`
		body := w.Body.String()
		if body != expected {
			t.Error(body, expected)
		}
	})
}

func TestReportDeclined(t *testing.T) {
	t.Run("declined_no_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		mockedReportHandler.ReportDeclined(w, r)
		status := w.Result().StatusCode
		if status != http.StatusUnauthorized {
			t.Error(status, http.StatusUnauthorized)
		}
	})

	t.Run("declined_with_auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.AddCookie(&http.Cookie{
			Name:     "access_token",
			Value:    "12345",
			HttpOnly: true,
		})
		r.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "-12345",
			HttpOnly: true,
		})
		nextHandler := http.HandlerFunc(mockedReportHandler.ReportDeclined)
		handlerToTest := mockedReportHandler.AuthMiddleware(nextHandler)

		handlerToTest.ServeHTTP(w, r)

		status := w.Result().StatusCode
		if status != http.StatusOK {
			t.Error(status, http.StatusOK)
		}
		expected := `{"type":"declined","content":"2222"}`
		body := w.Body.String()
		if body != expected {
			t.Error(body, expected)
		}
	})
}

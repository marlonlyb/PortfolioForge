package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/infrastructure/handlers"
)

type loginHandlerStub struct{}

func (loginHandlerStub) AdminLogin(c echo.Context) error        { return c.NoContent(http.StatusOK) }
func (loginHandlerStub) RequestEmailLogin(c echo.Context) error { return c.NoContent(http.StatusOK) }
func (loginHandlerStub) VerifyEmailLogin(c echo.Context) error  { return c.NoContent(http.StatusOK) }
func (loginHandlerStub) LoginWithGoogle(c echo.Context) error   { return c.NoContent(http.StatusOK) }

var _ handlers.LoginHandler = loginHandlerStub{}

func TestPublicAuthRoutesRejectLegacyEmailPasswordEndpoints(t *testing.T) {
	e := echo.New()
	handler := loginHandlerStub{}
	LoginPublic(e, handler)
	LoginAdmin(e, handler)

	legacyRequests := []string{
		"/api/v1/public/login",
		"/api/v1/public/register",
	}

	for _, path := range legacyRequests {
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(`{"email":"ada@example.com","password":"secret-123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("POST %s status = %d, want %d", path, rec.Code, http.StatusNotFound)
		}
	}
}

package cloudpocket

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCloudPocket(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	if assert.NoError(t, GetAllCloudPocket(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}

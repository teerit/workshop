package pocket

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCloudPocketNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"})
	mock.ExpectQuery("select * from pockets where id = \\$1").WillReturnRows(mockRows)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/cloud-pocket", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &handler{db: db}
	h.Delete(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, res.Code)
	}
}

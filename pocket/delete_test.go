//go:build unit

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
	h.DeleteCloudPocketById(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, res.Code)
	}
}

func TestDeleteCloudPocketBalanceNotEqualsZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
		AddRow(1, "Test Pocket", "Saving", "THB", 12)

	mock.ExpectQuery("select id, name, category, currency, balance from pockets where id = \\$1").WithArgs(sqlmock.AnyArg()).WillReturnRows(mockRows)
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/cloud-pocket", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &handler{db: db}
	h.DeleteCloudPocketById(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, res.Code)
	}
}

func TestDeleteCloudPocketSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
		AddRow(1, "Test Pocket", "Saving", "THB", 0)

	mockRowsToDeleted := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("select id, name, category, currency, balance from pockets where id = \\$1").WithArgs(sqlmock.AnyArg()).WillReturnRows(mockRows)
	mock.ExpectQuery("delete from pockets where id = \\$1 RETURNING id").WithArgs(sqlmock.AnyArg()).WillReturnRows(mockRowsToDeleted)
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/cloud-pocket", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &handler{db: db}
	h.DeleteCloudPocketById(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusAccepted, res.Code)
	}
}

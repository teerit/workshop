package pocket

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var dumpData = []Pocket{
	{
		ID:       12345,
		Name:     "Travel Fund",
		Category: "Vacation",
		Currency: "THB",
		Balance:  100,
	},
	{
		ID:       67890,
		Name:     "Savings",
		Category: "Emergency Fund",
		Currency: "THB",
		Balance:  200,
	},
}

func TestGetAllCloudPocket(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mockSql := "SELECT id, name, catagory, currency, balance FROM pockets"
	mockRow := sqlmock.NewRows([]string{"id", "name", "catagory", "currency", "balance"}).
		AddRow(dumpData[0].ID, dumpData[0].Name, dumpData[0].Category, dumpData[0].Currency, dumpData[0].Balance).
		AddRow(dumpData[1].ID, dumpData[1].Name, dumpData[1].Category, dumpData[1].Currency, dumpData[1].Balance)
	mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectQuery().WillReturnRows(mockRow)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	h := &handler{db: db}

	if assert.NoError(t, h.GetAllCloudPocket(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}

func TestGetCloudPocketByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mockSql := "SELECT id, name, catagory, currency, balance FROM pockets WHERE id = $1"
	mockRow := sqlmock.NewRows([]string{"id", "name", "catagory", "currency", "balance"}).
		AddRow(dumpData[0].ID, dumpData[0].Name, dumpData[0].Category, dumpData[0].Currency, dumpData[0].Balance)
	mock.ExpectPrepare(regexp.QuoteMeta(mockSql)).ExpectQuery().WillReturnRows(mockRow)

	cpkID := "12345"
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(cpkID)

	h := &handler{db: db}

	if assert.NoError(t, h.GetCloudPocketByID(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}

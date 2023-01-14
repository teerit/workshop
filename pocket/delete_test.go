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

func TestDeleteCloudPocketById(t *testing.T) {
	tests := []struct {
		name             string
		pathParam        string
		mockRows         *sqlmock.Rows
		mockRowToDeleted *sqlmock.Rows
		expectedStatus   int
	}{
		{
			name:             "TestDeleteCloudPocketNotFound",
			pathParam:        "1",
			mockRows:         sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}),
			mockRowToDeleted: sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedStatus:   http.StatusNotFound,
		},
		{
			name:      "TestDeleteCloudPocketBalanceNotEqualsZero",
			pathParam: "1",
			mockRows: sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
				AddRow(1, "Test Pocket", "Saving", "THB", 12),
			mockRowToDeleted: sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:      "TestDeleteCloudPocketSuccess",
			pathParam: "1",
			mockRows: sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
				AddRow(1, "Test Pocket", "Saving", "THB", 0),
			mockRowToDeleted: sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedStatus:   http.StatusAccepted,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			mock.ExpectQuery("select id, name, category, currency, balance from pockets where id = \\$1").WithArgs(sqlmock.AnyArg()).WillReturnRows(test.mockRows)
			mock.ExpectQuery("delete from pockets where id = \\$1 RETURNING id").WithArgs(sqlmock.AnyArg()).WillReturnRows(test.mockRowToDeleted)

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/cloud-pocket", nil)
			res := httptest.NewRecorder()

			c := e.NewContext(req, res)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.pathParam)
			h := &handler{db: db}
			h.DeleteCloudPocketById(c)

			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedStatus, res.Code)
			}
		})
	}
}

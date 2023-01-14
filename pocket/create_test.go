//go:build unit

package pocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePockets(t *testing.T) {

	tests := []struct {
		name       string
		cfgFlag    config.FeatureFlag
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{"create pockets succesfully",
			config.FeatureFlag{},
			func() (*sql.DB, error) {

				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				if err != nil {
					return nil, err
				}

				row := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(crestePkStmt).WithArgs("Travel Fund", "Vacation", "THB", 100.0).WillReturnRows(row)
				return db, err
			},
			`{"name":"Travel Fund","category":"Vacation","currency":"THB","balance":100.0}`,
			http.StatusCreated,
			`{"id": 1, "name":"Travel Fund","category":"Vacation","currency":"THB","balance":100.0}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := tc.sqlFn()
			h := New(tc.cfgFlag, db)
			// Assertions
			assert.NoError(t, err)
			if assert.NoError(t, h.CreatePocket(c)) {
				assert.Equal(t, tc.wantStatus, rec.Code)
				assert.JSONEq(t, tc.wantBody, rec.Body.String())
			}
		})
	}
}

func TestCreatePocket_Error(t *testing.T) {

	tests := []struct {
		name    string
		cfgFlag config.FeatureFlag
		sqlFn   func() (*sql.DB, error)
		reqBody string
		wantErr error
	}{
		{"create with bad request",
			config.FeatureFlag{},
			func() (*sql.DB, error) {
				return nil, nil
			},
			`{"name":Travel Fund,"category":"Vacation","currency":"THB","balance":100}`,
			echo.NewHTTPError(http.StatusBadRequest, "bad request body"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, _ := tc.sqlFn()
			h := New(tc.cfgFlag, db)

			berr := h.CreatePocket(c)
			// Assertions
			assert.Equal(t, berr, tc.wantErr)
		})
	}
}

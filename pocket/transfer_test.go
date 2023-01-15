//go:build unit

package pocket

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	sourcePocket = &pocket{
		ID:       1,
		Name:     "Travel Fund",
		Category: "Vacation",
		Currency: "THB",
		Balance:  1000.0,
	}
	destPocket = &pocket{
		ID:       2,
		Name:     "Savings",
		Category: "Emergency Fund",
		Currency: "THB",
		Balance:  1000.0,
	}
)

func TestTransferSuccess(t *testing.T) {
	// Arrange

	var transferTests = []struct {
		name     string
		cfgFlag  config.FeatureFlag
		sqlFn    func() (*sql.DB, error)
		payload  string
		wantCode int
		wantResp string
	}{
		{
			name:    "should return success transfer",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()

				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(destPocket.ID, destPocket.Name, destPocket.Category, destPocket.Currency, destPocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("2").WillReturnRows(newsMockRows)

				// mock update
				mockSqlUpdate := "UPDATE pockets"
				newsMockRows = sqlmock.NewRows([]string{"balance"}).
					AddRow(sourcePocket.Balance)
				mock.ExpectQuery(mockSqlUpdate).WithArgs(sourcePocket.ID, 950.0).
					WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"balance"}).
					AddRow(destPocket.Balance)
				mock.ExpectQuery(mockSqlUpdate).WithArgs(destPocket.ID, 1050.0).
					WillReturnRows(newsMockRows)

				// mock insert transaction
				newsMockRows = sqlmock.NewRows([]string{"id"}).
					AddRow("123")
				mock.ExpectQuery("INSERT INTO transactions (.+) RETURNING id").
					WithArgs("1", "2", 50.0, "Transfer from Travel fund to savings", "Success").
					WillReturnRows(newsMockRows)

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 50.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusOK,
			wantResp: `{
				"transaction_id": 123,
					"source_cloud_pocket": {
					"id": 1,
						"name": "Travel Fund",
						"category": "Vacation",
						"currency":"THB",
						"balance": 1000.00
				},
				"destination_cloud_pocket": {
					"id": 2,
						"name": "Savings",
						"category": "Emergency Fund",
						"currency":"THB",
						"balance": 1000.00
				},
				"status": "Success"
			}`,
		},
		{
			name:    "should bind fail return bad request",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				return nil, nil
			},
			payload: `{
				"source_cloud_pocket_id": 1234,
				"destination_cloud_pocket_id": "2",
				"amount": 50.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error_message": "Bad request",
				"status": "Failed"
			}`,
		},
		{
			name:    "should find source pocket fail return not found",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, _, _ := sqlmock.New()
				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 50.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusNotFound,
			wantResp: `{
				"error_message": "Not found source pocket",
				"status": "Failed"
			}`,
		},
		{
			name:    "should find destination pocket fail return not found",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 50.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusNotFound,
			wantResp: `{
				"error_message": "Not found destination pocket",
				"status": "Failed"
			}`,
		},
		{
			name:    "should return internal server error when pockets not enough balance and update transaction fail",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(destPocket.ID, destPocket.Name, destPocket.Category, destPocket.Currency, destPocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("2").WillReturnRows(newsMockRows)

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 1000000.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error_message": "Internal server error",
				"status": "Failed"
			}`,
		},
		{
			name:    "should return internal server error when pockets not enough balance and update transaction success",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(destPocket.ID, destPocket.Name, destPocket.Category, destPocket.Currency, destPocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("2").WillReturnRows(newsMockRows)

				// mock insert transaction
				newsMockRows = sqlmock.NewRows([]string{"id"}).
					AddRow("123")
				mock.ExpectQuery("INSERT INTO transactions (.+) RETURNING id").
					WithArgs("1", "2", 1000000.0, "Transfer from Travel fund to savings", "Failed").
					WillReturnRows(newsMockRows)

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 1000000.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error_message": "Not enough balance in the source cloud pocket",
				"status": "Failed"
			}`,
		},
		{
			name:    "should return internal server error when update source pocket fail",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(destPocket.ID, destPocket.Name, destPocket.Category, destPocket.Currency, destPocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("2").WillReturnRows(newsMockRows)

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 10.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error_message": "Internal server error",
				"status": "Failed"
			}`,
		},
		{
			name:    "should return internal server error when update destination pocket fail",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(destPocket.ID, destPocket.Name, destPocket.Category, destPocket.Currency, destPocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("2").WillReturnRows(newsMockRows)

				// mock update
				mockSqlUpdate := "UPDATE pockets"
				newsMockRows = sqlmock.NewRows([]string{"balance"}).
					AddRow(sourcePocket.Balance)
				mock.ExpectQuery(mockSqlUpdate).WithArgs(sourcePocket.ID, 950.0).
					WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"balance"}).
					AddRow(destPocket.Balance)
				mock.ExpectQuery(mockSqlUpdate).WithArgs(destPocket.ID, 1050.0).
					WillReturnError(errors.New("update pocket error"))

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 50.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error_message": "Internal server error",
				"status": "Failed"
			}`,
		},
		{
			name:    "should return internal server when insert transaction error",
			cfgFlag: config.FeatureFlag{},
			sqlFn: func() (*sql.DB, error) {
				db, mock, _ := sqlmock.New()
				// mock find
				mockSqlFind := "SELECT (.+) FROM pockets"
				newsMockRows := sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(sourcePocket.ID, sourcePocket.Name, sourcePocket.Category, sourcePocket.Currency, sourcePocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("1").WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(destPocket.ID, destPocket.Name, destPocket.Category, destPocket.Currency, destPocket.Balance)
				mock.ExpectQuery(mockSqlFind).WithArgs("2").WillReturnRows(newsMockRows)

				// mock update
				mockSqlUpdate := "UPDATE pockets"
				newsMockRows = sqlmock.NewRows([]string{"balance"}).
					AddRow(sourcePocket.Balance)
				mock.ExpectQuery(mockSqlUpdate).WithArgs(sourcePocket.ID, 950.0).
					WillReturnRows(newsMockRows)

				newsMockRows = sqlmock.NewRows([]string{"balance"}).
					AddRow(destPocket.Balance)
				mock.ExpectQuery(mockSqlUpdate).WithArgs(destPocket.ID, 1050.0).
					WillReturnRows(newsMockRows)

				return db, nil
			},
			payload: `{
				"source_cloud_pocket_id": "1",
				"destination_cloud_pocket_id": "2",
				"amount": 50.00,
				"description":"Transfer from Travel fund to savings"
			}`,
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error_message": "Internal server error",
				"status": "Failed"
			}`,
		},
	}

	for _, trt := range transferTests {
		t.Run(trt.name, func(t *testing.T) {
			// Arrange
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/cloud-pockets/transfer", strings.NewReader(trt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			db, _ := trt.sqlFn()
			h := &handler{db: db}

			wantStatus := trt.wantCode
			wantResp := trt.wantResp

			// Act
			gotErr := h.Transfer(c)
			gotResp := res.Body.String()
			gotStatus := res.Code

			// Assert
			assert.Nil(t, gotErr)
			assert.Equal(t, wantStatus, gotStatus)
			assert.JSONEq(t, wantResp, gotResp)
		})
	}
}

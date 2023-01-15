//go:build integration

package pocket

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDeletePocketsByIdIT(t *testing.T) {
	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}
	cfgFlag := config.FeatureFlag{}
	h := New(cfgFlag, sql)
	rec := CreatePocket(h)
	var p Pocket
	err = json.Unmarshal(rec.Body.Bytes(), &p)
	if err != nil {
		t.Error(err)
	}
	rec = h.DeletePocketByID(p)
	assert.Equal(t, http.StatusAccepted, rec.Code)
}
func (h *handler) DeletePocketByID(p Pocket) *httptest.ResponseRecorder {
	e := echo.New()
	e.DELETE("/cloud-pocket/:id", h.DeleteCloudPocketById)
	req := httptest.NewRequest(http.MethodDelete, uri("cloud-pocket", strconv.Itoa(int(p.ID))), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}
func CreatePocket(h *handler) *httptest.ResponseRecorder {
	e := echo.New()
	e.POST("/cloud-pockets", h.CreatePocket)
	reqBody := `{"name":"Travel Fund","category":"Vacation","currency":"THB","balance":0}`
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

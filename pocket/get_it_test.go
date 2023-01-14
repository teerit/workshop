//go:build integration

package pocket

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/kkgo-software-engineering/workshop/config"
)

func TestGetAllCloudPocketsIT(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	seedPocket(t)
	var cpk []Pocket
	res := request(t, http.MethodGet, uri("cloud-pockets"), nil)
	err := res.Decode(&cpk)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(cpk), 0)
}

func TestGetCloudPocketByIDIT(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	c := seedPocket(t)
	var cpk Pocket
	res := request(t, http.MethodGet, uri("cloud-pockets", strconv.Itoa(int(c.ID))), nil)
	err := res.Decode(&cpk)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func setup(t *testing.T) func() {
	e := echo.New()
	cfg := config.New().All()
	sql, _ := sql.Open("postgres", cfg.DBConnection)
	cfgFlag := config.FeatureFlag{}
	h := New(cfgFlag, sql)
	e.GET("/cloud-pockets", h.GetAllCloudPocket)
	e.GET("/cloud-pockets/:id", h.GetCloudPocketByID)
	e.POST("/cloud-pockets", h.CreatePocket)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)
	go func() {
		e.Start(addr)
	}()
	for {
		conn, _ := net.DialTimeout("tcp", fmt.Sprint("localhost:", os.Getenv("PORT")), 30*time.Second)
		if conn != nil {
			conn.Close()
			break
		}
	}

	teardown := func() {
		ctx, down := context.WithTimeout(context.Background(), 10*time.Second)
		defer down()
		err := e.Shutdown(ctx)
		assert.NoError(t, err)
	}

	return teardown
}

func seedPocket(t *testing.T) Pocket {
	var cpk Pocket
	body := bytes.NewBufferString(`{"name":"Travel Fund","category":"Vacation","currency":"THB","balance":100.0}`)

	err := request(t, http.MethodPost, uri("cloud-pockets"), body).Decode(&cpk)
	if err != nil {
		t.Fatal("can't create cloud pocket:", err)
	}
	return cpk
}

func uri(paths ...string) string {
	host := fmt.Sprint("http://localhost:", os.Getenv("PORT"))
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(t *testing.T, method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

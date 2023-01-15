//go:build integration

package pocket

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func TestTransferPocketsIt(t *testing.T) {
	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}

	cfgFlag := config.FeatureFlag{}
	h := New(cfgFlag, sql)
	e := echo.New()
	go func(e *echo.Echo) {
		e.POST("/cloud-pockets", h.CreatePocket)
		e.POST("/cloud-pockets/transfer", h.Transfer)
		e.Start(":2565")
	}(e)
	for {
		conn, err := net.DialTimeout("tcp", "localhost:2565", 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	createBody := `{"name":"Travel Fund","category":"Vacation","currency":"THB","balance":100.0}`
	trfBody := `{
		"source_cloud_pocket_id": 1,
		"destination_cloud_pocket_id": 1,
		"amount": 1.00,
		"description":"Transfer from Travel fund to savings"
	}`
	reSrc := regexp.MustCompile(`"source_cloud_pocket_id": \d+`)
	reDest := regexp.MustCompile(`"destination_cloud_pocket_id": \d+`)

	var srcPocket Pocket
	body := bytes.NewBufferString(createBody)
	res := Request(http.MethodPost, Uri("cloud-pockets"), body)
	err = res.Decode(&srcPocket)
	trfBody = reSrc.ReplaceAllString(trfBody, fmt.Sprintf(`"source_cloud_pocket_id": "%d"`, srcPocket.ID))

	var dstPocket Pocket
	body = bytes.NewBufferString(createBody)
	res = Request(http.MethodPost, Uri("cloud-pockets"), body)
	err = res.Decode(&dstPocket)
	trfBody = reDest.ReplaceAllString(trfBody, fmt.Sprintf(`"destination_cloud_pocket_id": "%d"`, dstPocket.ID))

	var trfResp transferResponse
	body = bytes.NewBufferString(trfBody)
	res = Request(http.MethodPost, Uri("cloud-pockets", "transfer"), body)
	err = res.Decode(&trfResp)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, trfResp.DestinationCloudPocket.Balance, 101.0)
	assert.Equal(t, trfResp.SourceCloudPocket.Balance, 99.0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = e.Shutdown(ctx)
}

func Uri(paths ...string) string {
	host := fmt.Sprintf("http://localhost%v", ":2565")
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func Request(method, url string, body io.Reader) *httpResponse {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &httpResponse{res, err}
}

type httpResponse struct {
	*http.Response
	err error
}

func (r *httpResponse) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

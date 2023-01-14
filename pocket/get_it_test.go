//go:build integration

package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCloudPocketsIT(t *testing.T) {

}

func setup(t *testing.T) func() {
	e := echo.New()
	// router.RegRoute(h)
	go func() {
		e.Start(os.Getenv("PORT"))
	}()
	for {
		conn, _ := net.DialTimeout("tcp", fmt.Sprint("localhost", os.Getenv("PORT")), 30*time.Second)
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

func seedExpense(t *testing.T) pocket {
	var cpk pocket
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	err := request(t, http.MethodPost, uri("cloud-pockets"), body).Decode(&cpk)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return cpk
}

func uri(paths ...string) string {
	host := fmt.Sprint("http://localhost", os.Getenv("PORT"))
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

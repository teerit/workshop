package pocket

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type transfer struct {
}

func (p *handler) Transfer(c echo.Context) error {

	return c.JSON(http.StatusOK, "hey")
}

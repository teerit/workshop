package pocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (p *handler) Transfer(c echo.Context) error {

	return c.JSON(http.StatusOK, "hey")
}

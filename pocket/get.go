package pocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllCloudPocket(c echo.Context) error {
	stmt, err := h.db.Prepare("SELECT id, name, catagory, currency, balance FROM pockets")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query all cloud pockets statment"})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't query all cloud pockets"})
	}

	pockets := []pocket{}
	for rows.Next() {
		var pck pocket
		err = rows.Scan(&pck.ID, &pck.Name, &pck.Category, &pck.Currency, &pck.Balance)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan cloud pocket"})
		}
		pockets = append(pockets, pck)
	}
	return c.JSON(http.StatusOK, pockets)
}

func GetCloudPocketByID(c echo.Context) error {
	// id := c.Param("id")
	// stmt, err := db.
	dumpData := pocket{
		ID:       12345,
		Name:     "Travel Fund",
		Category: "Vacation",
		Currency: "THB",
		Balance:  100,
	}

	return c.JSON(http.StatusOK, dumpData)
}

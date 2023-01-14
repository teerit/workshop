package pocket

import (
	"database/sql"
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
		var cpk pocket
		err = rows.Scan(&cpk.ID, &cpk.Name, &cpk.Category, &cpk.Currency, &cpk.Balance)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan cloud pocket"})
		}
		pockets = append(pockets, cpk)
	}
	return c.JSON(http.StatusOK, pockets)
}

func (h *handler) GetCloudPocketByID(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.db.Prepare("SELECT id, name, catagory, currency, balance FROM pockets WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query cloud pocket statment"})
	}
	row := stmt.QueryRow(id)
	cpk := pocket{}
	err = row.Scan(&cpk.ID, &cpk.Name, &cpk.Category, &cpk.Currency, &cpk.Balance)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNoContent, Err{Message: "Cloud Pocket Data Not Found"})
	case nil:
		return c.JSON(http.StatusOK, cpk)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan cloud pocket"})
	}
}

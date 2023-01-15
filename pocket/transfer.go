package pocket

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type transferDto struct {
	SourcePocketId string  `json:"source_cloud_pocket_id"`
	DestPocketId   string  `json:"destination_cloud_pocket_id"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
}

type transferResponse struct {
	TransactionId          int     `json:"transaction_id"`
	SourceCloudPocket      *Pocket `json:"source_cloud_pocket"`
	DestinationCloudPocket *Pocket `json:"destination_cloud_pocket"`
	Status                 string  `json:"status"`
}

func (h *handler) Transfer(c echo.Context) error {
	tDto := &transferDto{}
	sourcePocket := &Pocket{}
	destPocket := &Pocket{}

	// bind dto
	err := c.Bind(tDto)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp{
			Status:       "Failed",
			ErrorMessage: "Bad request",
		})
	}

	// query and bind source pocket
	err = findPocket(h.db, tDto.SourcePocketId, sourcePocket)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not found source pocket",
		})
	}
	// query and bind destination pocket
	err = findPocket(h.db, tDto.DestPocketId, destPocket)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not found destination pocket",
		})
	}

	if tDto.Amount > sourcePocket.Balance {
		_, err := insertTransaction(h.db, tDto, "Failed")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errorResp{
				Status:       "Failed",
				ErrorMessage: "Internal server error",
			})
		}
		return c.JSON(http.StatusBadRequest, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not enough balance in the source cloud pocket",
		})
	}

	// update amount source and destination
	begin, err := h.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}
	defer func() {
		switch err {
		case nil:
			err = begin.Commit()
		default:
			begin.Rollback()
		}
	}()

	err = updatePocket(begin, sourcePocket, sourcePocket.Balance-tDto.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	err = updatePocket(begin, destPocket, destPocket.Balance+tDto.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	tId, err := insertTransaction(h.db, tDto, "Success")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	resp := &transferResponse{
		TransactionId:          tId,
		SourceCloudPocket:      sourcePocket,
		DestinationCloudPocket: destPocket,
		Status:                 "Success",
	}

	return c.JSON(http.StatusOK, resp)
}

func findPocket(db *sql.DB, pid string, pk *Pocket) error {
	row := db.QueryRow(
		"SELECT * FROM pockets where id=$1",
		pid,
	)
	err := row.Scan(
		&pk.ID,
		&pk.Name,
		&pk.Category,
		&pk.Currency,
		&pk.Balance,
	)
	return err
}

func insertTransaction(db *sql.DB, tDto *transferDto, status string) (int, error) {
	row := db.QueryRow(
		"INSERT INTO transactions (source_pid, dest_pid, amount, description, date, status) values ($1, $2, $3, $4, current_timestamp, $5)  RETURNING id",
		tDto.SourcePocketId,
		tDto.DestPocketId,
		tDto.Amount,
		tDto.Description,
		status,
	)

	var resultId int
	err := row.Scan(&resultId)
	return resultId, err
}

func updatePocket(db *sql.Tx, p *Pocket, amount float64) error {
	row := db.QueryRow(
		"UPDATE pockets SET balance=$2 WHERE id=$1 RETURNING balance",
		p.ID,
		amount,
	)
	err := row.Scan(&p.Balance)
	return err
}

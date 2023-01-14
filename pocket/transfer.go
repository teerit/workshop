package pocket

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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

func (p *handler) Transfer(c echo.Context) error {
	model := &transferDto{}
	sourcePocket := &Pocket{}
	descPocket := &Pocket{}

	err := c.Bind(model)
	if err != nil {
		log.Println(err)

	}
	log.Printf("%v", model)

	// query source pocket
	err = findPocket(p.db, model.SourcePocketId, sourcePocket)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not found source pocket",
		})
	}
	// query destination pocket
	err = findPocket(p.db, model.DestPocketId, descPocket)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not found destination pocket",
		})
	}

	// check amount > source amount
	if model.Amount > sourcePocket.Balance {
		_, err := insertTransaction(p.db, model, "Failed")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errorResp{
				Status:       "Failed",
				ErrorMessage: "Internal server error",
			})
		}
		return c.JSON(http.StatusNotFound, errorResp{
			Status:       "Failed",
			ErrorMessage: "Not enough balance in the source cloud pocket",
		})
	}

	// update amount source and desc
	sourcePocket.Balance = sourcePocket.Balance - model.Amount
	descPocket.Balance = descPocket.Balance + model.Amount

	err = updatePocket(p.db, sourcePocket.ID, sourcePocket.Balance)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}
	err = updatePocket(p.db, descPocket.ID, descPocket.Balance)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	// insert transactions
	tId, err := insertTransaction(p.db, model, "Success")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp{
			Status:       "Failed",
			ErrorMessage: "Internal server error",
		})
	}

	resp := &transferResponse{
		TransactionId:          tId,
		SourceCloudPocket:      sourcePocket,
		DestinationCloudPocket: descPocket,
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
	if err != nil {
		return err
	}
	return nil
}

func insertTransaction(db *sql.DB, model *transferDto, status string) (int, error) {
	row := db.QueryRow(
		"INSERT INTO transactions (source_pid, dest_pid, amount, description, date, status) values ($1, $2, $3, $4, current_timestamp, $5)  RETURNING id",
		model.SourcePocketId,
		model.DestPocketId,
		model.Amount,
		model.Description,
		status,
	)

	var resultId int
	err := row.Scan(&resultId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return resultId, nil
}

func updatePocket(db *sql.DB, pid int64, result float64) error {
	stmt, err := db.Prepare("UPDATE pockets SET balance=$2 WHERE id=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pid, result)
	if err != nil {
		return err
	}

	return nil
}

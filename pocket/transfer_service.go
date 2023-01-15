package pocket

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type transferService struct {
	db  *sql.Tx
	log *zap.Logger
}

func newTransferService(db *sql.Tx, log *zap.Logger) transferService {
	return transferService{db: db, log: log}
}

func (t transferService) balanceCheck(tDto *transferDto, pocketBalance float64) (bool, error) {
	if tDto.Amount > pocketBalance {
		_, err := t.insertTransaction(tDto, "Failed")
		if err != nil {
			t.log.Error("insert transaction error")
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (t transferService) updatePocket(p *Pocket, amount float64) error {
	t.log.Info(fmt.Sprintf("transfer balance pocket id: %v", p.ID))
	row := t.db.QueryRow(
		"UPDATE pockets SET balance=$2 WHERE id=$1 RETURNING balance",
		p.ID,
		amount,
	)
	err := row.Scan(&p.Balance)
	return err
}

func (t transferService) findPocket(pid string, pk *Pocket) error {
	t.log.Info(fmt.Sprintf("find pocket id: %v", pk.ID))
	row := t.db.QueryRow(
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

func (t transferService) insertTransaction(tDto *transferDto, status string) (int, error) {
	t.log.Info("insert transaction from transfer")
	row := t.db.QueryRow(
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

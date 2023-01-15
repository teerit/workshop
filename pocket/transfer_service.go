package pocket

import "database/sql"

type transferService struct {
	db *sql.Tx
}

func newTransferService(db *sql.Tx) transferService {
	return transferService{db: db}
}

func (t transferService) balanceCheck(tDto *transferDto, pocketBalance float64) (bool, error) {
	if tDto.Amount > pocketBalance {
		_, err := t.insertTransaction(tDto, "Failed")
		if err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (t transferService) updatePocket(p *Pocket, amount float64) error {
	row := t.db.QueryRow(
		"UPDATE pockets SET balance=$2 WHERE id=$1 RETURNING balance",
		p.ID,
		amount,
	)
	err := row.Scan(&p.Balance)
	return err
}

func (t transferService) findPocket(pid string, pk *Pocket) error {
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

package pocket

import "database/sql"

type transferService struct {
	db *sql.DB
}

func newTransferService(db *sql.DB) transferService {
	return transferService{db: db}
}

func (t transferService) balanceCheck(tDto *transferDto, pocketBalance float64) (bool, error) {
	if tDto.Amount > pocketBalance {
		_, err := insertTransaction(t.db, tDto, "Failed")
		if err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

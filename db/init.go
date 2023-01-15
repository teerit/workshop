package db

import "database/sql"

func InitDatabase(db *sql.DB) error {
	createTb := `
			CREATE TABLE IF NOT EXISTS pockets(
				id SERIAL PRIMARY KEY,
				name TEXT,
				category TEXT,
				currency TEXT,
				balance float8
			);
`
	createTransaction := `
		CREATE TABLE IF NOT EXISTS transactions(
			id SERIAL PRIMARY KEY,
			source_pid INT,
			dest_pid INT,
			amount float8,
			description TEXT,
			date timestamp,
			status TEXT
		);
`
	_, err := db.Exec(createTb)
	if err != nil {
		return err
	}

	_, err = db.Exec(createTransaction)
	if err != nil {
		return err
	}
	return nil
}

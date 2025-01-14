package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func StartMigration(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "tuser" (
		id serial PRIMARY KEY NOT NULL,
		name VARCHAR(30) NOT NULL,
		password VARCHAR(30) NOT NULL, 
		role VARCHAR(30) NOT NULL,
		email   VARCHAR(30) NOT NULL,
		interests TEXT,
		refreshtkn TEXT,
		refreshttl INT
	);`)
	if err != nil {
		panic(err)
	}
	return err
}

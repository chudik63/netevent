package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func StartMigration(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "tuser" (
		id serial PRIMARY KEY NOT NULL,
		name VARCHAR(30),
		password VARCHAR(30), 
		role VARCHAR(30),
		email   VARCHAR(30),
		interest TEXT,
		accesstkn TEXT,
		accessttl INT,
		refreshtkn TEXT,
		refreshttl INT
	);`)
	if err != nil {
		panic(err)
	}
	return err
}

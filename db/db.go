package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

/*
create table init_table
(
    id         SERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL
);

create table update_table
(
    id         SERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL
);

create table terminate_table
(
    id         SERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL
);
*/

var db *sql.DB

func New() {
	connStr := "postgres://root:password@localhost:5432/load_tester?sslmode=disable"

	// Open a connection to the database.
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to open connection to DB: %v\n", err)
	}
}

func Close() {
	db.Close()
}

func InsertInitTable(accountID string, sessionID string) {
	_, err := db.Exec("INSERT INTO init_table (account_id, session_id) VALUES ($1, $2)", accountID, sessionID)
	if err != nil {
		fmt.Printf("failed to insert into table1: %v\n", err)
	}
}

func InsertUpdateTable(accountID string, sessionID string) {
	_, err := db.Exec("INSERT INTO update_table (account_id, session_id) VALUES ($1, $2)", accountID, sessionID)
	if err != nil {
		fmt.Printf("failed to insert into table2: %v\n", err)
	}
}

func InsertTerminateTable(accountID string, sessionID string) {
	_, err := db.Exec("INSERT INTO terminate_table (account_id, session_id) VALUES ($1, $2)", accountID, sessionID)
	if err != nil {
		fmt.Printf("failed to insert into table3: %v\n", err)
	}
}

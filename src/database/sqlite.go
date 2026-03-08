package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var DB *sql.DB

func usersTable(db *sql.DB) {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        username TEXT,
		username_md5 TEXT,
		password_md5 TEXT,
		xp INT DEFAULT 0,
		st INT DEFAULT 0,
		dx INT DEFAULT 0,
		iq INT DEFAULT 0,
		ht INT DEFAULT 0,
		points INT DEFAULT 0,
		credits INT DEFAULT 0,
		gold INT DEFAULT 0,
		ranking INT DEFAULT 0,
		totalRK INT DEFAULT 0,
		level INT DEFAULT 0,
		pmx INT DEFAULT 0,
		clanTag TEXT
    );
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Error().Msgf("Error creating users table: %s", err)
		return
	}

	log.Info().Msg("Users table created successfully")
}

func Initialize() {
	var err error
	DB, err = sql.Open("sqlite3", "./botzin.db")

	if err != nil {
		log.Error().Msgf("Error initializing SQLite database: %s", err)
		return
	}

	usersTable(DB)
}

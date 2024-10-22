package internal

import (
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path"
	"raccoonstash/internal/repository"
)

var DB *sql.DB
var Queries *repository.Queries

//go:embed sql/schema.sql
var databaseSchemaSql string

func InitializeDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", path.Join(Config.DataDir, "database.db"))
	if err != nil {
		println("Failed to open the database!")
		println(err.Error())
		return
	}

	// Initialize database if necessary
	{
		row := DB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='files'")
		var count int
		err := row.Scan(&count)
		if err != nil {
			log.Println("Error while ")
			log.Println(err)
		}
		if count == 0 {
			println("Table files does not exist, creating...")
			_, err := DB.Exec(databaseSchemaSql)
			if err != nil {
				println("Table creation failed!")
				println(err)
				return
			}
		}
	}

	Queries = repository.New(DB)
}

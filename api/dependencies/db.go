package dependencies

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type Database interface {
	Init(dataSourceName string) error
	Migrate() error
	GetConnection() *sqlx.DB
}

type PostgresSql struct {
	connection        *sqlx.DB
	connectionString  string
	migrationFileName string
}

func NewPostgreSql() *PostgresSql {
	return &PostgresSql{}
}

func (pg *PostgresSql) Init(dataSourceName string) error {
	pg.connectionString = dataSourceName
	pg.migrationFileName = "data.sql"

	var err error
	pg.connection, err = sqlx.Open("postgres", pg.connectionString)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return err
	}

	return nil
}

func (pg *PostgresSql) Migrate() error {
	db := pg.GetConnection()

	if db == nil {
		log.Fatalf("No database connection open, failed to migrate.")
	}

	raw, err := os.ReadFile(pg.migrationFileName)
	if err != nil {
		log.Printf("Error reading migration file, no migrations applied: %v", err)
		return err
	}

	if len(raw) == 0 {
		log.Println("Sql script empty. Nothing to be done.")
		return nil
	}

	tx := db.MustBegin()
	tx.MustExec(string(raw))
	tx.Commit()

	log.Println("Migrations complete")
	return nil
}

func (pg *PostgresSql) GetConnection() *sqlx.DB {
	return pg.connection
}

type NilDatabase struct{}

func (db *NilDatabase) Init(dataSourceName string) error { return nil }
func (db *NilDatabase) Migrate() error                   { return nil }
func (db *NilDatabase) GetConnection() *sqlx.DB           { return nil }

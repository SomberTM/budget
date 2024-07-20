package dependencies

import (
	"context"
	"database/sql"
	"log"
	"os"
)

type Database interface {
	Init(dataSourceName string) error
	Migrate() error
	GetConnection() *sql.DB
}

type PostgresSql struct {
	connection        *sql.DB
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
	pg.connection, err = sql.Open("postgres", pg.connectionString)
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
		log.Printf("Error reading sql: %v", err)
		return err
	}

	if len(raw) == 0 {
		log.Println("Sql script empty. Nothing to be done.")
		return nil
	}

	ctx, rollbackTx := context.WithCancel(context.Background())
	defer rollbackTx()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	_, err = db.Query(string(raw))
	if err != nil {
		log.Printf("Error running migration, rolling back: %v", err)
		return err
	}

	tx.Commit()
	log.Println("Migrations complete")
	return nil
}

func (pg *PostgresSql) GetConnection() *sql.DB {
	return pg.connection
}

type NilDatabase struct{}

func (db *NilDatabase) Init(dataSourceName string) error { return nil }
func (db *NilDatabase) Migrate() error                   { return nil }
func (db *NilDatabase) GetConnection() *sql.DB           { return nil }

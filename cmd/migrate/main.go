package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?multiStatements=true", user, pass, host, name)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Could not connect to MySQL: %v", err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatalf("Could not create mysql driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}

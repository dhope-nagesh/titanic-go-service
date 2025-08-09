package main

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	csvFile, err := os.Open("./data/titanic.csv")
	if err != nil {
		log.Fatalf("failed to open csv file 'titanic.csv': %v. Make sure it's in the root directory.", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Read() // Skip header

	db, err := sql.Open("sqlite3", "./data/titanic.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := `
	DROP TABLE IF EXISTS passengers;
	CREATE TABLE passengers (
		PassengerId INTEGER PRIMARY KEY,
		Survived INTEGER,
		Pclass INTEGER,
		Name TEXT,
		Sex TEXT,
		Age REAL,
		SibSp INTEGER,
		Parch INTEGER,
		Ticket TEXT,
		Fare REAL,
		Cabin TEXT,
		Embarked TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO passengers(PassengerId, Survived, Pclass, Name, Sex, Age, SibSp, 
                       Parch, Ticket, Fare, Cabin, Embarked) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatalf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading record: %v", err)
			continue
		}

		args := make([]interface{}, len(record))
		for i, val := range record {
			if val == "" {
				args[i] = nil
			} else {
				// Convert fields to their proper types
				if i == 5 || i == 9 { // Age or Fare
					args[i], _ = strconv.ParseFloat(val, 64)
				} else if i == 0 || i == 1 || i == 2 || i == 6 || i == 7 { // Integer fields
					args[i], _ = strconv.Atoi(val)
				} else { // String fields
					args[i] = val
				}
			}
		}

		_, err = stmt.Exec(args...)
		if err != nil {
			log.Printf("failed to insert record %+v: %v", record, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit transaction: %v", err)
	}

	log.Println("Database setup complete. titanic.db created successfully.")
}

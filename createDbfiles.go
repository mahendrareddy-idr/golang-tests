package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	outputDir    = `C:\Users\Mahendra\Documents\GO-LANG\Dataset\Databases&Structured Files`
	filesPerType = 100
)

func main() {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for i := 1; i <= filesPerType; i++ {
		createSQLite(fmt.Sprintf("data_%02d.db", i))
		createSQLite(fmt.Sprintf("data_%02d.sqlite", i))
		createSQL(fmt.Sprintf("data_%02d.sql", i))
		createJSON(fmt.Sprintf("data_%02d.json", i))
		createXML(fmt.Sprintf("data_%02d.xml", i))
	}

	fmt.Println("✅ Database & structured files created successfully")
}

// --------------------------------------------------
// SQLite (.db / .sqlite)

func createSQLite(name string) {
	path := filepath.Join(outputDir, name)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT,
			email TEXT
		);
	`)

	db.Exec(`
		INSERT INTO users (name, email)
		VALUES
		('Alice', 'alice@example.com'),
		('Bob', 'bob@example.com');
	`)
}

// --------------------------------------------------
// SQL dump (.sql)

func createSQL(name string) {
	content := `
CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	name TEXT,
	email TEXT
);

INSERT INTO users VALUES (1, 'Alice', 'alice@example.com');
INSERT INTO users VALUES (2, 'Bob', 'bob@example.com');
`
	writeFile(name, content)
}

// --------------------------------------------------
// JSON

func createJSON(name string) {
	content := `{
  "users": [
    { "id": 1, "name": "Alice", "email": "alice@example.com" },
    { "id": 2, "name": "Bob", "email": "bob@example.com" }
  ]
}`
	writeFile(name, content)
}

// --------------------------------------------------
// XML

func createXML(name string) {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<users>
  <user>
    <id>1</id>
    <name>Alice</name>
    <email>alice@example.com</email>
  </user>
  <user>
    <id>2</id>
    <name>Bob</name>
    <email>bob@example.com</email>
  </user>
</users>`
	writeFile(name, content)
}

// --------------------------------------------------

func writeFile(name, content string) {
	path := filepath.Join(outputDir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

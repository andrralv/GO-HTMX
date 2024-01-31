package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type manga struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./manga.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS manga (
			id TEXT PRIMARY KEY,
			title TEXT,
			author TEXT,
			quantity INTEGER
		);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	insertDataQuery := `
		INSERT INTO manga (id, title, author, quantity) VALUES
			('1', 'Monster', 'Naoki Urasawa', 3),
			('2', 'Full Metal Alchemist', 'Hiromu Arakawa', 3),
			('3', 'Slam Dunk', 'Takehiko Inoue', 3);
	`

	_, err = db.Exec(insertDataQuery)
	if err != nil {
		fmt.Println("Error inserting data:", err)
	}
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}

func getAllMangas() ([]manga, error) {
	rows, err := db.Query("SELECT id, title, author, quantity FROM manga")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangas []manga
	for rows.Next() {
		var m manga
		err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.Quantity)
		if err != nil {
			return nil, err
		}
		mangas = append(mangas, m)
	}

	return mangas, nil
}

func getMangaById(id string) (*manga, error) {
	var m manga
	err := db.QueryRow("SELECT id, title, author, quantity FROM manga WHERE id = ?", id).
		Scan(&m.ID, &m.Title, &m.Author, &m.Quantity)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func updateManga(m *manga) error {
	_, err := db.Exec("UPDATE manga SET title=?, author=?, quantity=? WHERE id=?", m.Title, m.Author, m.Quantity, m.ID)
	return err
}

func insertManga(m manga) error {
	_, err := db.Exec("INSERT INTO manga (id, title, author, quantity) VALUES (?, ?, ?, ?)", m.ID, m.Title, m.Author, m.Quantity)
	return err
}

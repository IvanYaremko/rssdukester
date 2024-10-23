package main

import (
	"database/sql"
	"log"

	"github.com/IvanYaremko/rssdukester/models/container"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalln("error opening sqlite3 database.db", err)
	}
	defer db.Close()

	if err := goose.Up(db, "./database.db"); err != nil {
		log.Fatalln("error migration database.db", err)
	}

	p := tea.NewProgram(container.Container{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

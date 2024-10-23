package main

import (
	"database/sql"
	"log"

	"github.com/IvanYaremko/rssdukester/models/container"
	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalln("error opening sqlite3 database.db", err)
	}
	defer db.Close()

	quries := database.New(db)

	p := tea.NewProgram(container.CreateContainer(quries), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

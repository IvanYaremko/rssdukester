package main

import (
	"database/sql"
	"log"

	"github.com/IvanYaremko/rssdukester/sql/database"
	"github.com/IvanYaremko/rssdukester/views"
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
	p := tea.NewProgram(views.InitHomeModel(quries), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

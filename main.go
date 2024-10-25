package main

import (
	"context"
	"database/sql"
	"fmt"
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

	feed, err := quries.GetFeedById(context.Background(), 0)
	if err != nil {
		log.Fatalln("error getting feed by id", err)
	}

	fmt.Println("feed:", feed)

	p := tea.NewProgram(container.CreateContainer(quries))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

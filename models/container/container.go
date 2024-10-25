package container

import (
	"github.com/IvanYaremko/rssdukester/config"
	"github.com/IvanYaremko/rssdukester/models/login"
	"github.com/IvanYaremko/rssdukester/models/rss"
	"github.com/IvanYaremko/rssdukester/sql/database"
	tea "github.com/charmbracelet/bubbletea"
)

type Container struct {
	dbQueries *database.Queries
}

func (c Container) Init() tea.Cmd {
	return userConfig
}

func (c Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case configFetched:
		if msg.cfg.Username == "" {
			model := login.InitialLoginModel()
			return model, model.Init()
		} else {
			model := rss.InitialiseFeedsModel(c.dbQueries)
			return model, model.Init()
		}
	case configFetchError:
		return c, func() tea.Msg {
			return configFetchError{err: msg.err}
		}
	}
	return c, nil
}

func (c Container) View() string {
	return ""
}

func CreateContainer(queries *database.Queries) Container {
	return Container{
		dbQueries: queries,
	}
}

func userConfig() tea.Msg {
	user, err := config.ReadConfig()
	if err != nil {
		return configFetchError{err}
	}
	return configFetched{
		cfg: user,
	}
}

type configFetched struct {
	cfg config.Config
}

type configFetchError struct {
	err error
}

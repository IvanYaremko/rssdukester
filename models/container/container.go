package container

import (
	"github.com/IvanYaremko/rssdukester/config"
	"github.com/IvanYaremko/rssdukester/models/login"
	"github.com/IvanYaremko/rssdukester/models/rss"
	tea "github.com/charmbracelet/bubbletea"
)

type Container struct {
}

func (c Container) Init() tea.Cmd {
	return userConfig
}

func (c Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case configFetched:
		if msg.cfg.Username == "" {
			return login.InitialLoginModel(), nil
		} else {
			return rss.InitialiseFeedsModel(), nil
		}
	case configFetchError:
		return c, func() tea.Msg {
			return configFetchError{err: msg.err}
		}
	}
	return c, nil
}

func (c Container) View() string {
	return "Hello"
}

func CreateContainer() Container {
	return Container{}
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

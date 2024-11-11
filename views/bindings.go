package views

import (
	"github.com/charmbracelet/bubbles/key"
)

var (
	helpBinding = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	quitBinding = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "exit"),
	)
	tabBinding = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next"),
	)
	strictUpBinding = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	)
	upBinding = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	)
	strictDownBinding = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	)
	downBinding = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	)
	shiftTabBinding = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous"),
	)
	enterBinding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "enter"),
	)
	escBinding = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
	backBinding = key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	)
	ctrlcBinding = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	addViewBiding = key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "Add view"),
	)
	removeBinding = key.NewBinding(
		key.WithKeys("X"),
		key.WithHelp("X", "remove"),
	)
)

type homeKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Add      key.Binding
	List     key.Binding
	Help     key.Binding
	Quit     key.Binding
	Ctrlc    key.Binding
}

func (k homeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit, k.Help}
}

func (k homeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Add},
		{k.List, k.Help, k.Quit},
		{k.Add},
	}
}

type addKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Help     key.Binding
	Back     key.Binding
	Ctrlc    key.Binding
}

func (a addKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{a.Enter, a.Back, a.Help, a.Ctrlc}
}

func (a addKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{a.Up, a.Down, a.Tab, a.ShiftTab},
		{a.Help, a.Back, a.Ctrlc},
	}
}

type rssItemKeyMap struct {
	Choose key.Binding
	Remove key.Binding
}

func (l rssItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		l.Choose,
		l.Remove,
	}
}

func (l rssItemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			l.Choose,
			l.Remove,
		},
	}
}

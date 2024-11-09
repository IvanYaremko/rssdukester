package views

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	ArrowUp   key.Binding
	Up        key.Binding
	ArrowDown key.Binding
	Down      key.Binding
	Tab       key.Binding
	ShiftTab  key.Binding
	Enter     key.Binding
	Add       key.Binding
	List      key.Binding
	Home      key.Binding
	Help      key.Binding
	Quit      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Add},
		{k.List, k.Help, k.Quit},
		{k.Home},
	}
}

var keys = keyMap{
	ArrowUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("'↑/k'", "up"),
	),
	ArrowDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("'↓/j'", "down"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("'↵'", "enter"),
	),
	Add: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("'A'", "add feed"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("'?'", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("'q'", "quit"),
	),
	Home: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("'H'", "home"),
	),
}

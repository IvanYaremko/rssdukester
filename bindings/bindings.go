package bindings

import "github.com/charmbracelet/bubbles/key"

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
	shiftTabBinding = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous"),
	)
	enterBinding = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "enter"),
	)
)

type HomeKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Add      key.Binding
	List     key.Binding
	Home     key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k HomeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit, k.Help}
}

func (k HomeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Add},
		{k.List, k.Help, k.Quit},
		{k.Home},
	}
}

var HomeKeys = HomeKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Tab:      tabBinding,
	ShiftTab: shiftTabBinding,
	Enter:    enterBinding,
	Add: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add feed"),
	),
	Help: helpBinding,
	Quit: quitBinding,
}

type AddKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Quit     key.Binding
	Home     key.Binding
	Help     key.Binding
}

func (a AddKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{a.Enter, a.Help, a.Quit}
}

func (a AddKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{a.Up, a.Down},
		{a.Tab, a.ShiftTab},
		{a.Help, a.Quit},
	}
}

var AddKeys = AddKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
	Tab:      tabBinding,
	ShiftTab: shiftTabBinding,
	Enter:    enterBinding,
	Help:     helpBinding,
	Quit:     quitBinding,
	Home: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "home page"),
	),
}

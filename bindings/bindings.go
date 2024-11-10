package bindings

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
)

type HomeKeyMap struct {
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

func (k HomeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit, k.Help}
}

func (k HomeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Add},
		{k.List, k.Help, k.Quit},
		{k.Add},
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
	Help:  helpBinding,
	Quit:  quitBinding,
	Ctrlc: ctrlcBinding,
}

type AddKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Quit     key.Binding
	Help     key.Binding
	Back     key.Binding
	Ctrlc    key.Binding
}

func (a AddKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{a.Enter, a.Back, a.Help, a.Quit}
}

func (a AddKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{a.Up, a.Down, a.Tab, a.ShiftTab},
		{a.Help, a.Back, a.Quit},
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
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Back:  escBinding,
	Ctrlc: ctrlcBinding,
}

type ListItemDelegateKeyMap struct {
	Choose key.Binding
	Remove key.Binding
}

func (l ListItemDelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		l.Choose,
		l.Remove,
	}
}

func (l ListItemDelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			l.Choose,
			l.Remove,
		},
	}
}

var ListItemDelegateKeys = ListItemDelegateKeyMap{
	Choose: enterBinding,
	Remove: key.NewBinding(
		key.WithKeys("x", "backspace"),
		key.WithHelp("x", "delete"),
	),
}

type ListKeysMap struct {
	Esc   key.Binding
	Back  key.Binding
	Ctrlc key.Binding
}

var ListKeys = ListKeysMap{
	Esc:   escBinding,
	Back:  backBinding,
	Ctrlc: ctrlcBinding,
}

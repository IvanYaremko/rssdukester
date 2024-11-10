package bindings

import "github.com/charmbracelet/bubbles/key"

var (
	helpBinding = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	quitBinding = key.NewBinding(
		key.WithKeys("Q"),
		key.WithHelp("Q", "exit"),
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
	Help     key.Binding
	Back     key.Binding
}

func (a AddKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{a.Enter, a.Back, a.Help, a.Quit}
}

func (a AddKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{a.Up, a.Down, a.Tab, a.ShiftTab},
		{a.Quit, a.Help, a.Back},
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
	Back:     escBinding,
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

func NewListItemDelegateKeyMap() ListItemDelegateKeyMap {
	return ListItemDelegateKeyMap{
		Choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		Remove: key.NewBinding(
			key.WithKeys("X", "backspace"),
			key.WithHelp("X", "delete"),
		),
	}
}

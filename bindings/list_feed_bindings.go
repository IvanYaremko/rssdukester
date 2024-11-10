package bindings

import "github.com/charmbracelet/bubbles/key"

type ListItemDelegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func (l ListItemDelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		l.choose,
		l.remove,
	}
}

func (l ListItemDelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			l.choose,
			l.remove,
		},
	}
}

func NewListItemDelegateKeyMap() ListItemDelegateKeyMap {
	return ListItemDelegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

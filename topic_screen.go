package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type topicScreen struct {
	model *model
}

func (t *topicScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return t.model, tea.Quit
		case tea.KeyEnter:
			t.model.selectedTopic = t.model.list.SelectedItem().FilterValue()
			currentScreen = &queryScreen{model: t.model}
			return t.model, textinput.Blink
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		t.model.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}
	var cmd tea.Cmd
	t.model.list, cmd = t.model.list.Update(msg)
	return t.model, cmd
}

func (t topicScreen) View() string {
	return docStyle.Render(t.model.list.View())
}

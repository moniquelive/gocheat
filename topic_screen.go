package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var str string
	if i, ok := listItem.(item); !ok {
		return
	} else {
		str = string(i)
	}
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}
	_, _ = fmt.Fprint(w, fn(str))
}

type topicScreen struct {
	model *model
}

func (t *topicScreen) Reset() {
	t.model.list.ResetFilter()
	t.model.queryTextInput.Reset()
}

func (t *topicScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return t.model, tea.Quit
		case tea.KeyEnter:
			t.model.selectedTopic = t.model.list.SelectedItem().FilterValue()
			setCurrentScreen(&queryScreen{model: t.model})
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

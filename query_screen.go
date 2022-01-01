package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type queryScreen struct {
	model *model
}

func (q *queryScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			q.model.selectedQuery = q.model.queryTextInput.Value()
			q.model.loadContents()
			q.model.viewport.SetContent(q.model.resultContent)
			currentScreen = &resultsScreen{model: q.model}
			return q.model, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return q.model, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		q.model.err = msg
		return q.model, nil
	}

	q.model.queryTextInput, cmd = q.model.queryTextInput.Update(msg)
	return q.model, cmd
}

func (q queryScreen) View() string {
	return fmt.Sprintf("Digite sua consulta\n\n%s\n\n(esc to quit)\n", q.model.queryTextInput.View())
}

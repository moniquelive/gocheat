package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type queryScreen struct {
	model *model
}

func (q *queryScreen) Reset() {
}

func (q *queryScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			q.model.selectedQuery = q.model.queryTextInput.Value()
			setCurrentScreen(&resultsScreen{model: q.model})
			return q.model, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return q.model, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		q.model.err = msg
		return q.model, nil
	}

	var cmd tea.Cmd
	q.model.queryTextInput, cmd = q.model.queryTextInput.Update(msg)
	return q.model, cmd
}

func (q queryScreen) View() string {
	return p.Sprintf("Type your query\n\n%s\n\n(esc to quit)\n", q.model.queryTextInput.View())
}

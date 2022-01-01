package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

const (
	headerHeight = 3
	footerHeight = 3
)

type resultsScreen struct {
	model *model
}

func (r *resultsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			currentScreen = &topicScreen{model: r.model}
			r.model.selectedTopic = ""
			r.model.selectedQuery = ""
			return r.model, nil
		case tea.KeyRunes:
			if msg.Runes[0] == 'q' {
				return r.model, tea.Quit
			}
		case tea.KeyCtrlC:
			return r.model, tea.Quit
		}

	case tea.WindowSizeMsg:
		verticalMargins := headerHeight + footerHeight
		r.model.viewport.Width = msg.Width
		r.model.viewport.Height = msg.Height - verticalMargins
	}

	// Because we're using the viewport's default update function (with pager-
	// style navigation) it's important that the viewport's update function:
	//
	// * Receives messages from the Bubble Tea runtime
	// * Returns commands to the Bubble Tea runtime
	//
	var cmd tea.Cmd
	r.model.viewport, cmd = r.model.viewport.Update(msg)
	return r.model, cmd
}

func (r resultsScreen) View() string {
	if !r.model.ready {
		return "\n  Initializing..."
	}

	// TODO: colocar a query no resultado
	headerTop := "╭───────────╮"
	headerMid := "│ Resultado ├"
	headerBot := "╰───────────╯"
	headerMid += strings.Repeat("─", r.model.viewport.Width-runewidth.StringWidth(headerMid))
	header := fmt.Sprintf("%s\n%s\n%s", headerTop, headerMid, headerBot)

	footerTop := "╭──────╮"
	footerMid := fmt.Sprintf("┤ %3.f%% │", r.model.viewport.ScrollPercent()*100)
	footerBot := "╰──────╯"
	gapSize := r.model.viewport.Width - runewidth.StringWidth(footerMid)
	footerTop = strings.Repeat(" ", gapSize) + footerTop
	footerMid = strings.Repeat("─", gapSize) + footerMid
	footerBot = strings.Repeat(" ", gapSize) + footerBot
	footer := fmt.Sprintf("%s\n%s\n%s", footerTop, footerMid, footerBot)

	return fmt.Sprintf("%s\n%s\n%s", header, r.model.viewport.View(), footer)
}

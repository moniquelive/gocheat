package main

import (
	"fmt"
	"io"
	"log"
	url2 "net/url"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/parnurzeal/gorequest"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type screen interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

var currentScreen screen

type model struct {
	err error

	list           list.Model
	queryTextInput textinput.Model
	viewport       viewport.Model

	ready         bool
	selectedTopic string
	selectedQuery string
	resultContent string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - (headerHeight + footerHeight)}
			m.viewport.YPosition = headerHeight + 1
		}
	}
	return currentScreen.Update(msg)
}

func (m model) View() string {
	return currentScreen.View()
}

func (m *model) loadContents() {
	query := m.selectedTopic
	if escapedQuery := url2.QueryEscape(m.selectedQuery); escapedQuery != "" {
		query += "/" + escapedQuery
	}
	url := fmt.Sprintf("https://cht.sh/%s?style=paraiso-dark", query)
	resp, _, errs := gorequest.New().
		Get(url).
		Set("User-Agent", "curl").
		End()
	if errs != nil {
		for _, err := range errs {
			log.Println("loadContents:", err)
		}
		currentScreen = &topicScreen{model: m}
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("main:", err)
	}
	m.resultContent = string(all)
}

func main() {
	m := InitialModel()
	currentScreen = &topicScreen{model: m}
	m.list.SetFilteringEnabled(true)
	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func InitialModel() *model {
	resp, _, errs := gorequest.New().
		Get("https://cht.sh/:list").
		Set("User-Agent", "curl").
		End()
	if errs != nil {
		for _, err := range errs {
			log.Println("main:", err)
		}
		os.Exit(-1)
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("main:", err)
	}

	var items []list.Item
	for _, each := range strings.Split(string(all), "\n") {
		items = append(items, item{title: each})
	}

	ti := textinput.NewModel()
	ti.Placeholder = "Digite sua consulta"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m := model{
		list:           list.NewModel(items, list.NewDefaultDelegate(), 0, 0),
		queryTextInput: ti,
	}
	m.list.Title = "Tópicos"
	return &m
}

package main

//go:generate gotext update -lang en-US,pt-BR -out catalog.go

import (
	"bytes"
	"fmt"
	"io"
	"log"
	url2 "net/url"
	"os"

	"golang.org/x/text/message"

	"github.com/Xuanwo/go-locale"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/parnurzeal/gorequest"
)

var p *message.Printer

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type screen interface {
	Reset()
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

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
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

func (m *model) loadContents() error {
	query := m.selectedTopic
	if m.selectedQuery != "" {
		if m.selectedQuery[0] == ':' {
			query += "/" + m.selectedQuery
		} else if escapedQuery := url2.QueryEscape(m.selectedQuery); escapedQuery != "" {
			query += "/" + escapedQuery
		}
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
		return errs[0]
	}
	defer resp.Body.Close()
	var all []byte
	var err error
	if all, err = io.ReadAll(resp.Body); err != nil {
		return err
	}
	m.resultContent = string(all)
	return nil
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

	lines := bytes.Split(all, []byte("\n"))
	var items = make([]list.Item, 0, len(lines))
	for _, each := range lines {
		items = append(items, item(each))
	}

	tm := textinput.NewModel()
	tm.Placeholder = p.Sprintf("Your query")
	tm.Focus()
	tm.CharLimit = 156
	tm.Width = 20

	lm := list.NewModel(items, itemDelegate{}, 0, 0)
	// TODO: enable filtering by default when it's available
	// https://github.com/charmbracelet/bubbles/issues/85
	m := model{
		list:           lm,
		queryTextInput: tm,
	}
	m.list.Title = p.Sprintf("Topics")
	return &m
}

func setCurrentScreen(scr screen) {
	currentScreen = scr
	currentScreen.Reset()
}

func main() {
	tag, err := locale.Detect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Tag:", tag)
	p = message.NewPrinter(tag)

	m := InitialModel()
	setCurrentScreen(&topicScreen{model: m})

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal("Error running program:", err)
	}
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

const (
	KeyColor    = "#A0D8EE"
	SepColor    = "#888888"
	DescColor   = "#A0D8EE"
	HeaderColor = "#A0D8EE"
)

const ncols = 3

type command struct {
	Desc    string
	Key     string
	Command string
}

type menu struct {
	Desc   string
	Key    string
	Action []command
	Menu   []menu
}

type config struct {
	KeyColor    string
	HeaderColor string
	SepColor    string
	DescColor   string
	Menu        []menu
	Action      []command
}

type model struct {
	Config      config
	CurrentMenu *menu
	err         error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keypress := msg.String()
		cMenu := *m.CurrentMenu
		for _, c := range cMenu.Menu {
			if c.Key == keypress {
				m.CurrentMenu = &c
				return m, nil
			}
		}
		for _, c := range cMenu.Action {
			if c.Key == keypress {
				fmt.Fprintln(os.Stdout, c.Command)
				return m, tea.Quit
			}
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func getKeyString(key, desc string, keyStyle, desStyle, parStyle lipgloss.Style) string {
	return strings.Join([]string{
		parStyle.Render("  ["),
		keyStyle.Render(key),
		parStyle.Render("]"),
		desStyle.Render(desc),
	}, "")
}

func (m model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	width, _, _ := term.GetSize(int(os.Stderr.Fd()))
	columnWidth := (width / 3) - 3
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.Config.HeaderColor))
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.Config.KeyColor))
	desStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.Config.DescColor)).
		PaddingLeft(2)
	parStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.Config.SepColor))
	padStyle := lipgloss.NewStyle().
		Width(columnWidth).
		Align(lipgloss.Left)

	c := *m.CurrentMenu
	i := 0
	var modeData = make([]string, ncols)
	var row string
	for _, k := range c.Menu {
		row = getKeyString(k.Key, k.Desc, keyStyle, desStyle, parStyle)
		modeData[i%ncols] += padStyle.Render(row) + "\n"
		i++
	}
	for _, k := range c.Action {
		row = getKeyString(k.Key, k.Desc, keyStyle, desStyle, parStyle)
		modeData[i%ncols] += padStyle.Render(row) + "\n"
		i++
	}

	return headerStyle.Render(c.Desc) + "\n" + lipgloss.JoinHorizontal(
		lipgloss.Top,
		modeData...,
	)
}

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		os.Exit(1)
	}

	locations := []string{
		os.ExpandEnv("$HOME/ghydra/config.toml"),
		configDir + "/ghydra/config.toml",
	}

	var configLocation string
	for _, configLocation = range locations {
		if _, err := os.Stat(configLocation); err == nil {
			break
		}
	}

	// Set default colors
	var cfg config
	cfg.KeyColor = KeyColor
	cfg.HeaderColor = HeaderColor
	cfg.SepColor = SepColor
	cfg.DescColor = DescColor

	// TODO: Handle config fallbacks better
	if _, err := toml.DecodeFile(configLocation, &cfg); err != nil {
		os.Exit(1)
	}

	Menu := menu{
		Desc:   "Main:",
		Key:    "",
		Action: cfg.Action,
		Menu:   cfg.Menu,
	}

	// Model initialization --> config and current menu setup
	m := model{}
	m.Config = cfg
	m.CurrentMenu = &Menu

	lipgloss.SetColorProfile(termenv.TrueColor)
	program := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

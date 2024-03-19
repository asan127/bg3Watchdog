package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type TickMsg time.Time
type LauncherDwellCompleted struct{}
type checkAppMsg struct{}

type model struct {
	spinner                spinner.Model
	appRunning             bool
	quitting               bool
	err                    error
	lastCheckTime          time.Time
	launcherDwellCompleted bool
	checkInterval          time.Duration
}

func (m model) Init() tea.Cmd {
	// Start the spinner
	spinnerCmd := m.spinner.Tick

	// Start initial check after LauncherDwell
	return tea.Batch(
		spinnerCmd,
		tea.Tick(LauncherDwell, func(t time.Time) tea.Msg {
			return LauncherDwellCompleted{}
		}),
	)
}

var count = 0

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	count = count + 1
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case LauncherDwellCompleted:
		m.launcherDwellCompleted = true
		m.checkInterval = PollInterval // Update checkInterval after initial dwell
		return m, tea.Tick(m.checkInterval, func(t time.Time) tea.Msg { return checkAppMsg{} })

	case checkAppMsg:
		if m.launcherDwellCompleted {
			// Check the application status only if the check interval has elapsed
			appRunning := isAppRunning()
			if appRunning != m.appRunning {
				m.appRunning = appRunning
			}
			m.lastCheckTime = time.Now()

			// If app is not running, quit
			if !m.appRunning {
				m.quitting = true
				return m, tea.Quit
			}
		}
		// Restart the ticker for the next check
		return m, tea.Tick(m.checkInterval, func(t time.Time) tea.Msg { return checkAppMsg{} })
	}

	var cmd tea.Cmd

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	status := "running"
	statusIcon := m.spinner.View()

	if !m.appRunning {
		status = "not running"
		statusIcon = `💩`
	}

	msg := "%s Target application 'bg3.exe' or 'bg3_dx11.exe' is %s.\n"

	if !Debug {
		return fmt.Sprintf(msg,
			statusIcon,
			status,
		)
	} else {
		msg += "Dwell completed: %t. Update cycle count: %d. Poll interval: %s, Time Since Last Check: %s"
		return fmt.Sprintf(msg,
			statusIcon,
			status,
			m.launcherDwellCompleted,
			count,
			m.checkInterval,
			time.Since(m.lastCheckTime),
		)
	}

}

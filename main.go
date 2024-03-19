package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ps "github.com/mitchellh/go-ps"
	"github.com/spf13/viper"
)

var (
	LauncherDwell time.Duration
	PollInterval  time.Duration
	Debug         bool
)

type AppConfig struct {
	SteamPath            string `mapstructure:"steam_path"`
	Bg3AppId             int    `mapstructure:"bg3_app_id"`
	LauncherDwellSeconds int    `mapstructure:"launcher_dwell_seconds"`
	MainAppPollSeconds   int    `mapstructure:"main_app_poll_seconds"`
	Debug                bool   `mapstructure:"debug"`
}

func main() {
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("steam_path", "C:\\Program Files (x86)\\Steam\\steam.exe")
	viper.SetDefault("bg3_app_id", 1086940)
	viper.SetDefault("launcher_dwell_seconds", 10)
	viper.SetDefault("main_app_poll_seconds", 5)
	viper.SetDefault("debug", false)

	// Read or create configuration
	if err := readOrCreateConfig(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Unmarshal configuration into struct
	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Error unmarshaling config:", err)
		os.Exit(2)
	}

	LauncherDwell = time.Duration(cfg.LauncherDwellSeconds) * time.Second
	PollInterval = time.Duration(cfg.MainAppPollSeconds) * time.Second
	Debug = cfg.Debug

	if Debug {
		fmt.Printf("Waiting %s after launching before monitoring...\n", LauncherDwell)
		fmt.Printf("Check every %s after launch dwell time for app closure...\n\n", PollInterval)
	}

	if !isAppRunning() {
		launchApp(cfg.SteamPath, cfg.Bg3AppId)
	}

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner:       s,
		appRunning:    true,
		lastCheckTime: time.Now(),
		checkInterval: LauncherDwell,
	}
}

func launchApp(SteamPath string, Bg3AppId int) {
	cmd := exec.Command(SteamPath, "-applaunch", strconv.Itoa(Bg3AppId))

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}

func isAppRunning() bool {
	processes, err := ps.Processes()
	if err != nil {
		return false
	}

	for _, proc := range processes {
		if strings.Contains(proc.Executable(), "bg3.exe") || strings.Contains(proc.Executable(), "bg3_dx11.exe") {
			return true
		}
	}

	return false
}

func readOrCreateConfig() error {
	// Check if the config file exists
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		// Write default configuration if file doesn't exist
		if err := viper.SafeWriteConfigAs("config.yaml"); err != nil {
			return err
		}
		fmt.Println("Default config file created: config.yaml")
	} else if err != nil {
		return err
	}

	// Read configuration from file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

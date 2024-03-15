package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/muesli/termenv"
)

const (
	checkInterval = 10 * time.Second // Check every 10 seconds
)

var (
	targetAppPath string
	targetAppName string
)

func main() {
	output := termenv.NewOutput(os.Stdout)

	if len(os.Args) < 2 {
		fmt.Println(output.String("Usage: ./sidecar <full_path_to_application>"))
		os.Exit(1)
	}

	targetAppPath = os.Args[1]
	targetAppName = filepath.Base(targetAppPath)

	// var launcherStarted = false

	// var bg3 ps.Process

	// if err := launchApp(targetAppPath); err != nil {
	// 	fmt.Printf("Failed to launch '%s': %v\n", targetAppName, err)
	// 	os.Exit(1)
	// }

	var cmd exec.Cmd

	for {

		if !isAppRunning() {
			fmt.Printf("Target application '%s' is not running, launching...\n", targetAppName)

			cmd := exec.Command(targetAppPath)
			err := cmd.Start()
			if err != nil {
				panic(err)
			}

			fmt.Printf("PID: %d\n", cmd.Process.Pid)

			// if err := launchApp(targetAppPath); err != nil {
			// 	fmt.Printf("Failed to launch '%s': %v\n", targetAppName, err)
			// 	os.Exit(1)
			// }
		}

		fmt.Println("Waiting for the thing to close...")
		err := cmd.Wait()

		log.Printf("Command finished with error: %v", err)

		fmt.Printf("Sleeping for %d seconds...", checkInterval)
		time.Sleep(checkInterval)
	}
}

// func main() {
//     for {
//         if !isAppRunning() {
//             fmt.Printf("Target application '%s' is not running, launching...\n", targetAppName)
//             if err := launchApp(); err != nil {
//                 fmt.Printf("Failed to launch '%s': %v\n", targetAppName, err)
//                 os.Exit(1)
//             }
//         }
//         time.Sleep(checkInterval)
//     }
// }

func isAppRunning() bool {
	return getProcess(targetAppName) != nil
}

func getProcess(targetAppName string) ps.Process {
	processes, err := ps.Processes()
	if err != nil {
		fmt.Printf("Error getting processes: %v\n", err)
	}

	for _, proc := range processes {
		if strings.Contains(proc.Executable(), targetAppName) {
			return proc
		}
	}

	return nil
}

func launchApp(targetAppPath string) error {
	cmd := exec.Command(targetAppPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

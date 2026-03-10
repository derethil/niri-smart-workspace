package main

import (
	"fmt"
	"io"
	"net"
	"niri-smart-workspace/internal/nirictl"
	"os"
)

func parseArgs() (isDebug, isDaemon bool, direction string) {
	for _, arg := range os.Args[1:] {
		if arg == "--debug" {
			isDebug = true
		} else if arg == "--daemon" {
			isDaemon = true
		} else if direction == "" {
			direction = arg
		}
	}

	if isDaemon {
		return isDebug, true, ""
	}

	if direction == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [--debug] <up|down>\n", os.Args[0])
		os.Exit(1)
	}

	if direction != "up" && direction != "down" {
		fmt.Fprintf(os.Stderr, "Invalid argument: %s (must be 'up' or 'down')\n", direction)
		os.Exit(1)
	}

	return isDebug, false, direction
}

func main() {
	isDebug, isDaemon, direction := parseArgs()

	if isDaemon {
		if err := nirictl.RunDaemon(isDebug); err != nil {
			fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	socketPath := nirictl.GetSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: daemon not running (start with --daemon)\n")
		os.Exit(1)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(direction)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if string(response) != "ok\n" {
		fmt.Fprintf(os.Stderr, "%s", string(response))
		os.Exit(1)
	}
}

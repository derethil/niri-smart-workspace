package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"niri-smart-workspace/internal/nirictl"
	"os"
)

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintln(os.Stderr)
	}
	os.Exit(1)
}

func parseArgs() (isDebug, isDaemon bool, direction string, err error) {
	flag.BoolVar(&isDebug, "debug", false, "enable debug logging")
	flag.BoolVar(&isDaemon, "daemon", false, "run as daemon")
	flag.Parse()

	if isDaemon {
		return isDebug, true, "", nil
	}

	if flag.NArg() < 1 {
		return false, false, "", fmt.Errorf("usage: %s [--debug] [--daemon] <up|down>", os.Args[0])
	}

	direction = flag.Arg(0)

	if direction != "up" && direction != "down" {
		return false, false, "", fmt.Errorf("invalid argument: %s (must be 'up' or 'down')", direction)
	}

	return isDebug, false, direction, nil
}

func main() {
	isDebug, isDaemon, direction, err := parseArgs()
	if err != nil {
		fatalf("%v", err)
	}

	if isDaemon {
		if err := nirictl.RunDaemon(isDebug); err != nil {
			fatalf("daemon error: %v", err)
		}
		return
	}

	socketPath := nirictl.GetSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fatalf("daemon not running (start with --daemon)")
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(direction)); err != nil {
		fatalf("failed to send command: %v", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		fatalf("failed to read response: %v", err)
	}

	if string(response) != "ok\n" {
		fatalf("%s", response)
	}
}

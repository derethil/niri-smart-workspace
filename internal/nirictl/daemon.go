package nirictl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func InitializeState() (*State, error) {
	var workspaces []Workspace
	var windows []Window

	debug("[INIT] Getting initial workspaces")

	cmd := exec.Command("niri", "msg", "--json", "workspaces")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(output, &workspaces)
	if err != nil {
		return nil, err
	}

	debug("[INIT] Getting initial windows")

	cmd = exec.Command("niri", "msg", "--json", "windows")
	output, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(output, &windows)
	if err != nil {
		return nil, err
	}

	return NewState(workspaces, windows), nil
}

func startEventListener(state *State) {
	go func() {
		debug("[EVENT] Starting event listener")
		cmd := exec.Command("niri", "msg", "--json", "event-stream")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get event stream: %v\n", err)
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start event stream: %v\n", err)
			return
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			var event Event
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				debug("[EVENT] Failed to parse event: %v", err)
				continue
			}

			if event.WorkspacesChanged != nil {
				debug("[EVENT] Workspaces changed: %d workspaces", len(event.WorkspacesChanged.Workspaces))
				state.UpdateWorkspaces(event.WorkspacesChanged.Workspaces)
			}
			if event.WindowsChanged != nil {
				debug("[EVENT] Windows changed: %d windows", len(event.WindowsChanged.Windows))
				state.UpdateWindows(event.WindowsChanged.Windows)
			}
			if event.WorkspaceActivated != nil {
				debug("[EVENT] Workspace activated: id=%d, focused=%v", event.WorkspaceActivated.ID, event.WorkspaceActivated.Focused)
				state.UpdateFocusedWorkspace(event.WorkspaceActivated.ID)
			}
			if event.WindowOpenedOrChanged != nil {
				debug("[EVENT] Window opened/changed: id=%d workspace=%d", event.WindowOpenedOrChanged.Window.ID, event.WindowOpenedOrChanged.Window.WorkspaceID)
				state.AddOrUpdateWindow(event.WindowOpenedOrChanged.Window)
			}
			if event.WindowClosed != nil {
				debug("[EVENT] Window closed: id=%d", event.WindowClosed.ID)
				state.RemoveWindow(event.WindowClosed.ID)
			}
		}
	}()
}

func handleConnection(conn net.Conn, state *State) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	direction := string(buf[:n])
	debug("[CONN] Received request: %s", direction)

	workspaces, windows := state.Get()
	debug("[CONN] Current state: %d workspaces, %d windows", len(workspaces), len(windows))

	if err := navigate(direction, workspaces, windows); err != nil {
		debug("[CONN] Navigation error: %v", err)
		fmt.Fprintf(conn, "error: %v\n", err)
	} else {
		debug("[CONN] Navigation successful")
		fmt.Fprintf(conn, "ok\n")
	}
}

func runSocketServer(state *State) error {
	socketPath := GetSocketPath()
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}
	defer os.Remove(socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, state)
	}
}

func RunDaemon(isDebug bool) error {
	SetDebugMode(isDebug)

	state, err := InitializeState()
	if err != nil {
		return err
	}

	startEventListener(state)

	return runSocketServer(state)
}

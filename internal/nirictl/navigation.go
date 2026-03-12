package nirictl

import (
	"fmt"
	"os/exec"
	"sort"
)

func countWindowsOnWorkspace(workspaceID int, windows []Window) int {
	count := 0
	for _, w := range windows {
		if w.WorkspaceID == workspaceID {
			count++
		}
	}
	return count
}

func findFocusedWorkspace(workspaces []Workspace) (*Workspace, error) {
	for i := range workspaces {
		if workspaces[i].IsFocused {
			return &workspaces[i], nil
		}
	}
	return nil, fmt.Errorf("no focused workspace found")
}

func filterWorkspacesByOutput(workspaces []Workspace, output string) []Workspace {
	var filtered []Workspace
	for _, ws := range workspaces {
		if ws.Output == output {
			filtered = append(filtered, ws)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Idx < filtered[j].Idx
	})

	return filtered
}

func buildWorkspacesWithWindowsMap(workspaces []Workspace, windows []Window) map[int]bool {
	workspacesWithWindows := make(map[int]bool)
	debug("[MAP] Building workspace map with %d windows:", len(windows))
	for _, w := range windows {
		debug("[MAP]   window id=%d workspace_id=%d", w.ID, w.WorkspaceID)
	}
	for _, ws := range workspaces {
		count := countWindowsOnWorkspace(ws.ID, windows)
		debug("[MAP] Workspace id=%d idx=%d has %d windows", ws.ID, ws.Idx, count)
		if count > 0 {
			workspacesWithWindows[ws.Idx] = true
		}
	}
	return workspacesWithWindows
}

func workspaceExists(workspaces []Workspace, idx int) bool {
	for _, ws := range workspaces {
		if ws.Idx == idx {
			return true
		}
	}
	return false
}

func handleUpWorkspace(currentIdx int, outputWorkspaces []Workspace, workspacesWithWindows map[int]bool) error {
	debug("[UP] currentIdx=%d, workspacesWithWindows=%v", currentIdx, workspacesWithWindows)
	for i := currentIdx - 1; i >= 1; i-- {
		debug("[UP] checking idx=%d, exists=%v, hasWindows=%v", i, workspaceExists(outputWorkspaces, i), workspacesWithWindows[i])
		if workspaceExists(outputWorkspaces, i) && workspacesWithWindows[i] {
			debug("[UP] navigating to workspace %d", i)
			return focusWorkspace(i)
		}
	}
	debug("[UP] no workspace found, doing nothing")
	return nil
}

func handleDownWorkspace(currentIdx int, outputWorkspaces []Workspace, workspacesWithWindows map[int]bool) error {
	debug("[DOWN] currentIdx=%d, workspacesWithWindows=%v", currentIdx, workspacesWithWindows)
	for i := currentIdx + 1; i <= len(outputWorkspaces)+10; i++ {
		debug("[DOWN] checking idx=%d, hasWindows=%v", i, workspacesWithWindows[i])
		if workspacesWithWindows[i] {
			debug("[DOWN] navigating to workspace %d", i)
			return focusWorkspace(i)
		}
	}
	debug("[DOWN] no workspace found, doing nothing")
	return nil
}

func navigate(direction string, workspaces []Workspace, windows []Window) error {
	debug("[NAVIGATE] direction=%s, workspaces=%d, windows=%d", direction, len(workspaces), len(windows))

	currentWorkspace, err := findFocusedWorkspace(workspaces)
	if err != nil {
		return err
	}
	debug("[NAVIGATE] current workspace: idx=%d, id=%d", currentWorkspace.Idx, currentWorkspace.ID)

	outputWorkspaces := filterWorkspacesByOutput(workspaces, currentWorkspace.Output)
	debug("[NAVIGATE] output workspaces: %d", len(outputWorkspaces))

	workspacesWithWindows := buildWorkspacesWithWindowsMap(outputWorkspaces, windows)

	if direction == "up" {
		return handleUpWorkspace(currentWorkspace.Idx, outputWorkspaces, workspacesWithWindows)
	}
	return handleDownWorkspace(currentWorkspace.Idx, outputWorkspaces, workspacesWithWindows)
}

func focusWorkspace(idx int) error {
	cmd := exec.Command("niri", "msg", "action", "focus-workspace", fmt.Sprintf("%d", idx))
	return cmd.Run()
}

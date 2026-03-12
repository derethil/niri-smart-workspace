package nirictl

import "sync"

type Workspace struct {
	ID             int     `json:"id"`
	Idx            int     `json:"idx"`
	Name           *string `json:"name"`
	Output         string  `json:"output"`
	IsActive       bool    `json:"is_active"`
	IsFocused      bool    `json:"is_focused"`
	ActiveWindowID *int    `json:"active_window_id"`
}

type Window struct {
	ID          int `json:"id"`
	WorkspaceID int `json:"workspace_id"`
}

type WorkspacesChangedEvent struct {
	Workspaces []Workspace `json:"workspaces"`
}

type WindowsChangedEvent struct {
	Windows []Window `json:"windows"`
}

type WorkspaceActivatedEvent struct {
	ID      int  `json:"id"`
	Focused bool `json:"focused"`
}

type WindowOpenedOrChangedEvent struct {
	Window Window `json:"window"`
}

type WindowClosedEvent struct {
	ID int `json:"id"`
}

type Event struct {
	WorkspacesChanged      *WorkspacesChangedEvent      `json:"WorkspacesChanged"`
	WindowsChanged         *WindowsChangedEvent         `json:"WindowsChanged"`
	WorkspaceActivated     *WorkspaceActivatedEvent     `json:"WorkspaceActivated"`
	WindowOpenedOrChanged  *WindowOpenedOrChangedEvent  `json:"WindowOpenedOrChanged"`
	WindowClosed           *WindowClosedEvent           `json:"WindowClosed"`
}

type State struct {
	mu         sync.RWMutex
	workspaces []Workspace
	windows    []Window
}

func NewState(workspaces []Workspace, windows []Window) *State {
	debug("[INIT] Loaded %d workspaces", len(workspaces))
	debug("[INIT] Loaded %d windows", len(windows))

	state := &State{
		workspaces: workspaces,
		windows:    windows,
	}

	for _, ws := range workspaces {
		if ws.IsFocused {
			state.UpdateFocusedWorkspace(ws.ID)
			debug("[INIT] Loaded focused workspace: id=%d, output=%s", ws.ID, ws.Output)
		}
	}

	return state
}

func (s *State) Get() ([]Workspace, []Window) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.workspaces, s.windows
}

func (s *State) UpdateWorkspaces(workspaces []Workspace) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workspaces = workspaces
}

func (s *State) UpdateWindows(windows []Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows = windows
}

func (s *State) UpdateFocusedWorkspace(workspaceID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.workspaces {
		s.workspaces[i].IsFocused = s.workspaces[i].ID == workspaceID
		s.workspaces[i].IsActive = s.workspaces[i].ID == workspaceID
	}
}

func (s *State) AddOrUpdateWindow(window Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.windows {
		if s.windows[i].ID == window.ID {
			s.windows[i] = window
			return
		}
	}
	s.windows = append(s.windows, window)
}

func (s *State) RemoveWindow(windowID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.windows {
		if s.windows[i].ID == windowID {
			s.windows = append(s.windows[:i], s.windows[i+1:]...)
			return
		}
	}
}

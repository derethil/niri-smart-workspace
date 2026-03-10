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

type Event struct {
	WorkspacesChanged  *WorkspacesChangedEvent  `json:"WorkspacesChanged"`
	WindowsChanged     *WindowsChangedEvent     `json:"WindowsChanged"`
	WorkspaceActivated *WorkspaceActivatedEvent `json:"WorkspaceActivated"`
}

type State struct {
	mu         sync.RWMutex
	isDebug    bool
	workspaces []Workspace
	windows    []Window
}

func NewState(isDebug bool, workspaces []Workspace, windows []Window) *State {
	debug("[INIT] Loaded %d workspaces", len(workspaces))
	debug("[INIT] Loaded %d windows", len(windows))

	state := &State{
		isDebug:    isDebug,
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

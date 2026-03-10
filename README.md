# niri-smart-workspace

Smart workspace navigation for the [niri](https://github.com/YaLTeR/niri)
Wayland compositor that automatically skips empty workspaces.

## Features

- Navigate to the next/previous workspace with windows, skipping empty ones
- Runs as a background daemon listening to niri events to avoid fetching state
  on every invocation
- Fast Unix socket IPC for navigation commands
- Lightweight and written in Go

## How It Works

The daemon subscribes to niri's event stream to maintain real-time state of all
workspaces and windows. When you trigger navigation (up/down), it finds the
nearest non-empty workspace in that direction and focuses it.

Navigation behavior:

- **up** (previous): Finds the previous workspace with windows, stops at
  workspace 1
- **down** (next): Finds the next workspace with windows, stops at the last
  workspace with windows

## Installation

### With Nix Flakes

Add to your configuration:

```nix
programs.niri.enable = true;

# The daemon will start automatically with your graphical session
```

Then bind keys to the navigation commands, e.g.:

```nix
"Mod+BracketLeft".action = spawn-sh "${getExe pkgs.niri-smart-workspace} up";
"Mod+BracketRight".action = spawn-sh "${getExe pkgs.niri-smart-workspace} down";
```

### From Source

```bash
cd src
go build -o niri-smart-workspace
```

## Usage

Start the daemon:

```bash
niri-smart-workspace --daemon
```

Navigate workspaces:

```bash
niri-smart-workspace up    
niri-smart-workspace down
```

Debug mode:

```bash
niri-smart-workspace --daemon --debug
```

## Architecture

- **Daemon**: Listens to `niri msg event-stream` and maintains workspace/window
  state
- **Client**: Sends navigation commands via Unix socket to
  `/run/user/$UID/niri-smart-workspace.sock`
- **Events**: Handles `WorkspacesChanged`, `WindowsChanged`, and
  `WorkspaceActivated` events from niri

## Requirements

- Go 1.25+
- niri compositor
- Linux with systemd (for daemon management)

## License

MIT

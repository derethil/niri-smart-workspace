package nirictl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var debugMode bool

func debug(format string, args ...any) {
	if debugMode {
		log.Printf(format, args...)
	}
}

func GetSocketPath() string {
	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if runtime == "" {
		runtime = "/tmp"
	}
	return filepath.Join(runtime, "niri-smart-workspace.sock")
}

func focusWorkspace(idx int) error {
	cmd := exec.Command("niri", "msg", "action", "focus-workspace", fmt.Sprintf("%d", idx))
	return cmd.Run()
}

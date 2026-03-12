package nirictl

import (
	"log"
	"os"
	"path/filepath"
)

var debugMode bool

func SetDebugMode(enabled bool) {
	debugMode = enabled
}

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

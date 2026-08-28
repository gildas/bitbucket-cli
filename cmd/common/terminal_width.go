package common

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/term"
)

var (
	terminalWidth     int
	terminalWidthOnce sync.Once
)

func GetTerminalWidth() int {
	terminalWidthOnce.Do(func() {
		terminalWidth = 80
		if os.Getenv("TMUX") != "" {
			if output, err := exec.Command("tmux", "display-message", "-p", "#{pane_width}").Output(); err == nil {
				if width, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil && width > 0 {
					terminalWidth = width
				}
			}
		}
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
			terminalWidth = width
		}
		terminalWidth = 80
	})
	return terminalWidth
}

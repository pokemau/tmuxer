package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Window struct {
	Name     string
	Layout   string
	Commands []string
}

type Config struct {
	Name    string
	Root    string
	Window []Window
}

func (c Config) createSession() {
	if err := exec.Command("tmux", "has-session", "-t", c.Name).Run(); err == nil {
		c.attachSession()
		return
	}

	if len(c.Window) == 0 {
		log.Fatal("No windows were provided")
	}

	initCmd := exec.Command("tmux", "new-session", "-d", "-s", c.Name)
	err := initCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	for _, window := range c.Window {
		c.processWindow(window)
	}

	// kill auto-created first window
	windows := c.getSessionWindows()
	if len(windows) != 0 {
		c.killWindow(windows[0])
		c.selectWindow(windows[0])
	}
}

func (c Config) processWindow(window Window) {
	c.createWindow(window)

	if len(window.Layout) != 0 {
		c.createPane(window)
	}

	panes := c.getPanes(window)

	if len(window.Commands) > 0 {
		for i, command := range window.Commands {
			if command != "" {
				c.sendKeys(window.Name, panes[i], command)
			}
		}
	}

	c.selectPane(window, panes[0])
}

func (c Config) createWindow(window Window) {
	// tmux new-window -t mysess -n server
	args := []string{"new-window", "-t", c.Name}

	if window.Name != "" {
		args = append(args, "-n", window.Name)
	}

	cmd := exec.Command("tmux", args...)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) sendKeys(windowName string, paneIndex int, command string) {
	// tmux send-keys -t my-session:editor.0 'nvim .' Enter
	target := c.Name + ":" + windowName + "." + strconv.Itoa(paneIndex)

	cmd := exec.Command("tmux", "send-keys", "-t", target, command, "Enter")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) selectPane(window Window, paneIndex int) {
	// tmux select-pane -t sess:win.0
	target := c.Name + ":" + window.Name + "." + strconv.Itoa(paneIndex)
	if err := exec.Command("tmux", "select-pane", "-t", target).Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) selectWindow(windowIndex int) {
	target := c.Name + ":" + strconv.Itoa(windowIndex)
	if err := exec.Command("tmux", "select-window", "-t", target).Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) killWindow(windowIndex int) {
	target := c.Name + ":" + strconv.Itoa(windowIndex)
	if err := exec.Command("tmux", "kill-window", "-t", target).Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) getSessionWindows() []int {
	out, err := exec.Command("tmux", "list-windows", "-t", c.Name, "-F", "#{window_index}").Output()
	if err != nil {
		log.Fatal(err)
	}

	var indices []int
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if i, err := strconv.Atoi(line); err == nil {
			indices = append(indices, i)
		}
	}
	return indices
}

func (c Config) getPanes(window Window) []int {
	target := c.Name + ":" + window.Name

	out, err := exec.Command("tmux", "list-panes", "-t", target, "-F", "#{pane_index}").Output()
	if err != nil {
		log.Fatal(err)
	}

	var indices []int
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if i, err := strconv.Atoi(line); err == nil {
			indices = append(indices, i)
		}
	}
	return indices
}

func (c Config) createPane(window Window) {
	// tmux split-window -h -t my-session:editor
	split := "-" + string(window.Layout[0])
	cmd := exec.Command("tmux", "split-window", split, "-t", c.Name+":"+window.Name)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) attachSession() {
	cmd := exec.Command("tmux", "attach", "-t", c.Name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

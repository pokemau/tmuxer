package main

import (
	"fmt"
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
	Windows []Window
}

func (c Config) createSession() {
	if err := exec.Command("tmux", "has-session", "-t", c.Name).Run(); err == nil {
		c.attachSession()
		return
	}

	if len(c.Windows) == 0 {
		fmt.Println("NO WINDOWS")
		return
	}

	initCmd := exec.Command("tmux", "new-session", "-d", "-s", c.Name)
	err := initCmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	for _, window := range c.Windows {
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

		panes := c.getPanes(window)

		if len(panes) != 0 {
			for i, paneIndex := range panes {
				c.sendKeys(window, &paneIndex, i)
			}
			c.selectPane(window, panes[0])
		}
	} else {
		c.sendKeys(window, nil, 0)
	}
}

func (c Config) createWindow(window Window) {
	// tmux new-window -t mysess -n server
	cmd := exec.Command("tmux", "new-window", "-t", c.Name, "-n", window.Name)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func (c Config) sendKeys(window Window, paneIndex *int, commandIndex int) {
	// tmux send-keys -t my-session:editor.0 'nvim .' Enter
	target := c.Name + ":" + window.Name
	if paneIndex != nil {
		target += "." + strconv.Itoa(*paneIndex)
	}
	cmd := exec.Command("tmux", "send-keys", "-t", target, window.Commands[commandIndex], "Enter")
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

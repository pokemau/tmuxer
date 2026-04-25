package main

import (
	"flag"
	"fmt"
	"log"

	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func main() {
	flag.Parse()

	args := flag.Args()

	configFile := "tm.toml"

	if len(args) > 0 {
		if args[0] == "init" {
			if len(args) < 2 {
				log.Fatal("Usage: tmuxer init <name> [config-path]")
			}

			name := args[1]
			configPath := ""

			if len(args) >= 3 {
				configPath = args[2]
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				configPath = cwd
			}

			generateConfig(name, configPath)
			return
		}

		info, err := os.Stat(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {
			configFile = filepath.Join(args[0], "tm.toml")
		} else {
			configFile = args[0]
		}
	}

	config := loadConfig(configFile)

	var conf Config

	_, err := toml.Decode(config, &conf)
	if err != nil {
		log.Fatal(err)
	}

	if len(conf.Name) == 0 {
		log.Fatal("Insert a session name")
	}

	conf.createSession()
	conf.attachSession()
}

func generateConfig(sessionName string, configPath string) {
	content := fmt.Sprintf(`# tmux session name
name = "%s"
root = "%s"

[[window]]
# name = ""
# layout = "horizontal" or "vertical"
commands = []
`, sessionName, configPath)

	outPath := filepath.Join(configPath, "tm.toml")

	if _, err := os.Stat(outPath); err == nil {
		log.Fatalf("%s already exists", outPath)
	}

	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created %s\n", outPath)
}

func loadConfig(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

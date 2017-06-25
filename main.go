package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

type BranchConfig struct {
	NoForcePush        bool `json:"noForcePush"`
	NoPush             bool `json:"noPush"`
	AskBeforeForcePush bool `json:"askBeforeForcePush"`
	AskBeforePush      bool `json:"askBeforePush"`
}

type Config map[string]BranchConfig

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Usage = "Push a branch"
	app.UsageText = "gopush <branch>"
	app.Author = "Bao Pham"
	app.Email = "gbaopham@gmail.com"
	app.EnableBashCompletion = true

	app.Before = func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			return errors.New("Branch is required")
		}

		return nil
	}

	app.BashComplete = bashComplete

	app.Action = push

	app.Run(os.Args)
}

func push(c *cli.Context) {
	var force bool
	var remote, branch string

	for _, arg := range c.Args() {
		if strings.HasPrefix(arg, "-") {
			force = arg == "-f" || arg == "--force"
		} else if remote == "" {
			remote = arg
		} else if branch == "" {
			branch = arg
		}
	}

	if branch == "" {
		branch, remote = remote, "origin"
	}

	config, err := getConfig()

	if err != nil {
		color.Red(err.Error())
		return
	}

	branchConfig, ok := config[branch]

	if !ok {
		execCommand(remote, branch, force)
		return
	}

	if force && branchConfig.NoForcePush {
		proceed := false
		if branchConfig.AskBeforeForcePush {
			proceed = prompt("Are you sure you want to force push to: %s?", branch)
			if !proceed || err != nil {
				return
			}
		}

		if !proceed {
			color.Red("Cannot force push to: %s", branch)
			return
		}
	}

	if !force && branchConfig.NoPush {
		proceed := false
		if branchConfig.AskBeforePush {
			proceed = prompt("Are you sure you want to push to: %s?", branch)
			if !proceed || err != nil {
				return
			}
		}

		if !proceed {
			color.Red("Cannot push directly to: %s", branch)
			return
		}
	}

	execCommand(remote, branch, force)
}

func bashComplete(c *cli.Context) {
	if c.NArg() > 0 {
		return
	}

	out, err := exec.Command("git", "branch").Output()
	sout := string(out)

	if err != nil || strings.Contains(sout, "Not a git repository") {
		return
	}

	branches := strings.Split(sout, "\n")

	for _, branch := range branches {
		branch = strings.TrimLeft(branch, "*")
		fmt.Println(strings.TrimSpace(branch))
	}
}

func execCommand(remote, branch string, force bool) {
	var cmd *exec.Cmd

	if force {
		color.Green("Force pushing to %s:%s %s...", remote, branch, "-f")
		cmd = exec.Command("git", "push", remote, branch, "-f")
	} else {
		color.Green("Pushing to %s:%s...", remote, branch)
		cmd = exec.Command("git", "push", remote, branch)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func getConfig() (Config, error) {
	user, err := user.Current()

	if err != nil {
		return Config{}, err
	}

	dir, err := os.Getwd()
	stop := false

	for !stop {
		file := path.Join(dir, ".gopush.json")

		if _, err := os.Stat(file); !os.IsNotExist(err) {
			out, err := ioutil.ReadFile(file)

			var config Config
			err = json.Unmarshal(out, &config)
			return config, err
		}

		if dir == user.HomeDir {
			stop = true
		}

		dir = path.Join(dir, "..")
	}

	return Config{}, nil
}

func prompt(format string, a ...interface{}) bool {
	if !strings.Contains(format, "[y/n]") {
		format += " [y/n] "
	}
	if len(a) == 0 {
		fmt.Print(format)
	} else {
		fmt.Printf(format, a...)
	}
	return handlePromptResponse()
}

func handlePromptResponse() bool {
	var response string
	_, err := fmt.Scanln(&response)

	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		return true
	}

	return false
}

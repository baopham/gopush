package main

import (
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

	app.BashComplete = func(c *cli.Context) {
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

	app.Action = func(c *cli.Context) {
		remote, branch := c.Args().Get(0), c.Args().Get(1)

		if branch == "" {
			branch = remote
			remote = "origin"
		}

		restrictedBranches, err := getRestrictedBranches()

		if err != nil {
			color.Red(err.Error())
			return
		}

		for _, restrictedBranch := range restrictedBranches {
			if restrictedBranch == branch {
				color.Red("Cannot push directly to: %s", branch)
				return
			}
		}

		cmd := exec.Command("git", "push", remote, branch)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	app.Run(os.Args)
}

func getRestrictedBranches() ([]string, error) {
	branches := []string{}
	user, err := user.Current()

	if err != nil {
		return branches, nil
	}

	dir, err := os.Getwd()

	if err != nil {
		return branches, nil
	}

	stop := false

	for !stop {
		file := dir + "/.gopush_restricted"

		if _, err := os.Stat(file); !os.IsNotExist(err) {
			config, err := ioutil.ReadFile(file)

			if err != nil {
				return branches, err
			}

			return strings.Split(string(config), "\n"), nil
		}

		if dir == user.HomeDir {
			stop = true
		}

		dir = path.Join(dir, "..")
	}

	return branches, nil
}

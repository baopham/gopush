package main

import (
	"errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Usage = "Push a branch"
	app.UsageText = "gopush <branch>"
	app.Author = "Bao Pham"
	app.Email = "gbaopham@gmail.com"

	app.Before = func(c *cli.Context) error {
		if len(c.Args()) < 1 {
			return errors.New("Branch is required")
		}

		return nil
	}

	app.Action = func(c *cli.Context) {
		remote, branch := c.Args().Get(0), c.Args().Get(1)

		if branch == "" {
			branch = remote
			remote = "origin"
		}

		restrictedBranches := getRestrictedBranches()

		for _, restrictedBranch := range restrictedBranches {
			if restrictedBranch == branch {
				log.Fatalf("Cannot push directly to: %s", branch)
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

func getRestrictedBranches() []string {
	user, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	path := user.HomeDir + "/.gopush"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []string{}
	}

	config, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(config), "\n")
}

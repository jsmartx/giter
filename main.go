package main

import (
	"log"
	"os"

	"github.com/jsmartx/giter/cmd"
	"github.com/urfave/cli"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Usage = "Git users manager"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all the git user config",
			Action:  cmd.List,
		},
		{
			Name:   "use",
			Usage:  "Change git user config to username",
			Action: cmd.Use,
		},
		{
			Name:   "show",
			Usage:  "Show git user detail",
			Action: cmd.Show,
		},
		{
			Name:    "add",
			Aliases: []string{"new"},
			Usage:   "Add one custom user config",
			Action:  cmd.Add,
		},
		{
			Name:   "update",
			Usage:  "Update one custom user config",
			Action: cmd.Update,
		},
		{
			Name:    "del",
			Aliases: []string{"rm"},
			Usage:   "Delete one custom user config",
			Action:  cmd.Delete,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

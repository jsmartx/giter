package cmd

import (
	"fmt"

	"github.com/jsmartx/giter/git"
	"github.com/jsmartx/giter/store"
	"github.com/urfave/cli"
)

func List(c *cli.Context) error {
	g, err := git.New(".")
	var cur *git.User
	if err == nil {
		cur = g.GetUser()
	}
	filter := c.Args().First()
	s := store.New()
	users := s.List(filter, false)
	for _, u := range users {
		if cur != nil && cur.Name == u.Name && cur.Email == u.Email {
			fmt.Printf(" * %s\n", u.String())
		} else {
			fmt.Printf("   %s\n", u.String())
		}
	}
	return nil
}

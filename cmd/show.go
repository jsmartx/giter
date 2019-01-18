package cmd

import (
	"errors"
	"fmt"

	"github.com/jsmartx/giter/git"
	"github.com/jsmartx/giter/store"
	"github.com/jsmartx/giter/util"
	"github.com/urfave/cli"
)

func Show(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		g, err := git.New(".")
		if err == nil {
			if u := g.GetUser(); u != nil {
				name = u.Name
			}
		}
		if name == "" {
			name = util.Prompt(&util.PromptConfig{
				Prompt: "user name: ",
			})
		}
	}
	s := store.New()
	users := s.List(name, true)
	if len(users) == 0 {
		return errors.New("user not found")
	}
	for i, u := range users {
		keyPath, err := u.KeyPath()
		if err != nil {
			return err
		}
		if i == 0 {
			fmt.Println()
		}
		if u.IsSSH() {
			fmt.Printf("       User: %s\n", u.String())
			fmt.Printf("      Email: %s\n", u.Email)
			fmt.Printf("Private Key: %s\n", keyPath)
			fmt.Printf(" Public Key: %s.pub\n", keyPath)
		} else {
			fmt.Printf("      User: %s\n", u.String())
			fmt.Printf("     Email: %s\n", u.Email)
			fmt.Printf("Credential: %s.credential\n", keyPath)
		}
		fmt.Println()
	}
	return nil
}

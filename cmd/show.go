package cmd

import (
	"errors"
	"fmt"
	"github.com/jsmartx/giter/git"
	"github.com/jsmartx/giter/store"
	"github.com/jsmartx/giter/util"
	"github.com/urfave/cli"
	"strconv"
)

func Show(c *cli.Context) error {
	g, err := git.New(".")
	if err != nil {
		return err
	}
	name := c.Args().First()
	if name == "" {
		defaultName := ""
		if u := g.GetUser(); u != nil {
			defaultName = u.Name
		}
		name = util.Prompt(&util.PromptConfig{
			Prompt:  "user name: ",
			Default: defaultName,
		})
	}
	s := store.New()
	users := s.List(name, true)
	if len(users) == 0 {
		return errors.New("User not found!")
	}
	u := users[0]
	if len(users) > 1 {
		fmt.Printf("There are %d users:\n", len(users))
		for i, item := range users {
			fmt.Printf("%4d) %s\n", i+1, item.String())
		}
		str := util.Prompt(&util.PromptConfig{
			Prompt: "Enter number to select user: ",
		})
		i, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if i < 1 || i > len(users) {
			return errors.New("Out of range")
		}
		u = users[i-1]
	}
	keyPath, err := u.KeyPath()
	if err != nil {
		return err
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
	return nil
}

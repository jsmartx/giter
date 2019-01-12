package cmd

import (
	"errors"
	"fmt"
	"github.com/jsmartx/giter/git"
	"github.com/jsmartx/giter/ssh"
	"github.com/jsmartx/giter/store"
	"github.com/jsmartx/giter/util"
	"github.com/urfave/cli"
	"strconv"
)

func Use(c *cli.Context) error {
	g, err := git.New(".")
	if err != nil {
		return err
	}
	name := c.Args().First()
	if name == "" {
		name = util.Prompt(&util.PromptConfig{
			Prompt: "user name: ",
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
	if u.IsSSH() {
		keyPath, err := u.KeyPath()
		if err != nil {
			fmt.Println(err)
		}
		s := ssh.New()
		err = s.SetHost(&ssh.Host{
			Key:          u.Host,
			Hostname:     u.Host,
			Port:         u.Port,
			IdentityFile: keyPath,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return g.SetUser(u.Name, u.Email)
}

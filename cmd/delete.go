package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jsmartx/giter/store"
	"github.com/jsmartx/giter/util"
	"github.com/urfave/cli"
)

func Delete(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		name = util.Prompt(&util.PromptConfig{
			Prompt: "user name: ",
		})
	}
	s := store.New()
	users := s.List(name, true)
	if len(users) == 0 {
		return errors.New("user not found")
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
	txt := util.Prompt(&util.PromptConfig{
		Prompt: fmt.Sprintf("Are you sure to delete '%s' [Y/n]? ", u.String()),
	})
	if txt != "y" && txt != "Y" {
		return nil
	}
	return s.Delete(u)
}

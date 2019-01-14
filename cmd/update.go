package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jsmartx/giter/store"
	"github.com/jsmartx/giter/util"
	"github.com/urfave/cli"
)

func Update(c *cli.Context) error {
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
	email := util.Prompt(&util.PromptConfig{
		Prompt:  "user email: ",
		Default: u.Email,
	})
	urlStr := util.Prompt(&util.PromptConfig{
		Prompt:  "git server: ",
		Default: u.FullHost(),
	})

	url, err := util.ParseURL(urlStr)
	util.CheckError(err)

	host, port := util.SplitHostPort(url.Host)
	user := &store.User{
		Name:   name,
		Email:  email,
		Scheme: url.Scheme,
		Host:   host,
		Port:   port,
	}
	options := &store.Options{}
	if user.IsSSH() {
		if u.IsSSH() {
			txt := util.Prompt(&util.PromptConfig{
				Prompt: fmt.Sprintf("Regenerate the SSH key [Y/n]? "),
			})
			if txt != "y" && txt != "Y" {
				p, err := u.KeyPath()
				if err != nil {
					return err
				}
				options.KeyPath = p
			}
		}
	} else {
		if !u.IsSSH() {
			txt := util.Prompt(&util.PromptConfig{
				Prompt: fmt.Sprintf("Reset password [Y/n]? "),
			})
			if txt == "y" || txt == "Y" {
				pwd := util.Prompt(&util.PromptConfig{
					Prompt: "user password: ",
					Silent: true,
				})
				options.Password = pwd
			}
		} else {
			pwd := util.Prompt(&util.PromptConfig{
				Prompt: "user password: ",
				Silent: true,
			})
			options.Password = pwd
		}
	}
	return s.Update(u.Hash(), user, options)
}

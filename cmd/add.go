package cmd

import (
	"fmt"

	"github.com/jsmartx/giter/git"
	"github.com/jsmartx/giter/store"
	"github.com/jsmartx/giter/util"
	"github.com/urfave/cli"
)

func getUser() *git.User {
	g, err := git.New(".")
	if err != nil {
		return git.GlobalUser()
	}
	return g.GetUser()
}

func getHost() string {
	g, err := git.New(".")
	if err != nil {
		return "ssh://github.com"
	}
	urls, err := g.Remotes()
	if err != nil || len(urls) == 0 {
		return "ssh://github.com"
	}
	return fmt.Sprintf("%s://%s", urls[0].Scheme, urls[0].Host)
}

func Add(c *cli.Context) error {
	u := getUser()
	nameCfg := &util.PromptConfig{
		Prompt: "user name: ",
	}
	emailCfg := &util.PromptConfig{
		Prompt: "user email: ",
	}
	urlCfg := &util.PromptConfig{
		Prompt:  "git server: ",
		Default: getHost(),
	}
	if u != nil {
		nameCfg.Default = u.Name
		emailCfg.Default = u.Email
	}
	name := util.Prompt(nameCfg)
	email := util.Prompt(emailCfg)
	urlStr := util.Prompt(urlCfg)

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
		pvtPath, err := util.SysSSHConfig()
		if err == nil {
			fmt.Printf("There is a SSH key: %s\nYou can use this key or generate a new SSH key.\n", pvtPath)
			txt := util.Prompt(&util.PromptConfig{
				Prompt: fmt.Sprintf("Do you want to use '%s' [Y/n]? ", pvtPath),
			})
			if txt == "y" || txt == "Y" {
				options.KeyPath = pvtPath
			}
		}
	} else {
		pwd := util.Prompt(&util.PromptConfig{
			Prompt: "user password: ",
			Silent: true,
		})
		options.Password = pwd
	}
	s := store.New()
	return s.Add(user, options)
}

package git

import (
	fs "io/ioutil"
	"net/url"

	"github.com/jsmartx/giter/util"
	homedir "github.com/mitchellh/go-homedir"
	git "gopkg.in/src-d/go-git.v4"
	config "gopkg.in/src-d/go-git.v4/config"
)

type Git struct {
	r *git.Repository
}

type User struct {
	Name  string
	Email string
}

func New(path string) (*Git, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	g := &Git{r: r}
	return g, nil
}

func GlobalUser() *User {
	p, err := homedir.Expand("~/.gitconfig")
	if err != nil {
		return nil
	}

	data, err := fs.ReadFile(p)
	if err != nil {
		return nil
	}

	cfg := config.NewConfig()
	if err = cfg.Unmarshal(data); err != nil {
		return nil
	}

	u := cfg.Raw.Section("user")
	name := u.Option("name")
	email := u.Option("email")
	if name == "" || email == "" {
		return nil
	}
	return &User{Name: name, Email: email}
}

func (g *Git) Remotes() ([]*url.URL, error) {
	cfg, err := g.r.Config()

	if err != nil {
		return nil, err
	}

	urls := make([]*url.URL, 0)

	for _, remote := range cfg.Remotes {
		for _, repo := range remote.URLs {
			repoURL, err := util.ParseURL(repo)
			if err != nil {
				continue
			}
			if remote.Name == "origin" {
				urls = append([]*url.URL{repoURL}, urls...)
			} else {
				urls = append(urls, repoURL)
			}
		}
	}
	return urls, nil
}

func (g *Git) GetUser() *User {
	cfg, err := g.r.Config()
	if err != nil {
		return nil
	}

	u := cfg.Raw.Section("user")
	name := u.Option("name")
	email := u.Option("email")
	if name == "" || email == "" {
		return nil
	}
	return &User{Name: name, Email: email}
}

func (g *Git) SetUser(name, email string) error {
	cfg, err := g.r.Config()

	if err != nil {
		return err
	}

	u := cfg.Raw.Section("user")
	u.SetOption("name", name)
	u.SetOption("email", email)
	return g.r.Storer.SetConfig(cfg)
}

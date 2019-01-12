package store

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsmartx/giter/util"
	fs "io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const ROOT = "~/.giter/"

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   string `json:"port"`
}

type Options struct {
	KeyPath  string
	Password string
}

func (u *User) Hash() string {
	host := util.JoinHostPort(u.Host, u.Port)
	url := fmt.Sprintf("%s://%s@%s", u.Scheme, u.Name, host)
	hash := md5.New()
	hash.Write([]byte(url))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (u *User) Url(pwd string) string {
	fullUrl := &url.URL{
		Scheme: u.Scheme,
		User:   url.UserPassword(u.Name, pwd),
		Host:   util.JoinHostPort(u.Host, u.Port),
	}
	return fullUrl.String()
}

func (u *User) KeyPath() (string, error) {
	return util.JoinPath(ROOT, "keys", u.Hash())
}

func (u *User) FullHost() string {
	host := util.JoinHostPort(u.Host, u.Port)
	return fmt.Sprintf("%s://%s", u.Scheme, host)
}

func (u *User) IsSSH() bool {
	if u.Scheme == "http" || u.Scheme == "https" {
		return false
	}
	return true
}

func (u *User) String() string {
	return fmt.Sprintf("%s - %s", u.Name, u.FullHost())
}

type Config struct {
	Users []*User `json:"users"`
}

type Store struct {
	c *Config
}

func loadConfig() *Config {
	cfgPath, err := util.JoinPath(ROOT, "config.json")
	if err != nil {
		return nil
	}
	f, err := os.Open(cfgPath)
	defer f.Close()
	if err != nil {
		return nil
	}
	var cfg Config
	parser := json.NewDecoder(f)
	parser.Decode(&cfg)
	return &cfg
}

func saveConfig(cfg *Config) error {
	root, err := util.Mkdir(ROOT)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(root, "config.json")
	// write marshaled data to the file
	return fs.WriteFile(cfgPath, b, 0755)
}

func New() *Store {
	cfg := loadConfig()
	if cfg != nil {
		return &Store{c: cfg}
	}
	cfg = &Config{
		Users: make([]*User, 0),
	}
	saveConfig(cfg)
	return &Store{c: cfg}
}

func (s *Store) check(user *User) error {
	for _, u := range s.c.Users {
		if u.Hash() == user.Hash() {
			return errors.New("User already exist!")
		}
	}
	return nil
}

func (s *Store) Add(u *User, opts *Options) error {
	if err := s.check(u); err != nil {
		return err
	}
	keysDir, err := util.Mkdir(ROOT, "keys")
	if err != nil {
		return err
	}
	if u.IsSSH() {
		pvtPath := filepath.Join(keysDir, u.Hash())
		pubPath := filepath.Join(keysDir, u.Hash()+".pub")
		if opts.KeyPath == "" {
			if err := util.Keygen(pubPath, pvtPath, 4096); err != nil {
				return err
			}
		} else {
			if err := util.Copy(opts.KeyPath, pvtPath); err != nil {
				return err
			}
			if err := util.Copy(opts.KeyPath+".pub", pubPath); err != nil {
				return err
			}
		}
	} else {
		pwdPath := filepath.Join(keysDir, u.Hash()+".credential")
		data := []byte(u.Url(opts.Password))
		if err := fs.WriteFile(pwdPath, data, 0755); err != nil {
			return err
		}
	}
	s.c.Users = append(s.c.Users, u)
	return saveConfig(s.c)
}

func (s *Store) Update(hash string, u *User, opts *Options) error {
	if hash != u.Hash() {
		if err := s.check(u); err != nil {
			return err
		}
	}
	keysDir, err := util.Mkdir(ROOT, "keys")
	if err != nil {
		return err
	}
	if u.IsSSH() {
		pvtPath := filepath.Join(keysDir, u.Hash())
		pubPath := filepath.Join(keysDir, u.Hash()+".pub")
		if opts.KeyPath == "" {
			if err := util.Keygen(pubPath, pvtPath, 4096); err != nil {
				return err
			}
		} else if opts.KeyPath != pvtPath {
			if err := util.Copy(opts.KeyPath, pvtPath); err != nil {
				return err
			}
			if err := util.Copy(opts.KeyPath+".pub", pubPath); err != nil {
				return err
			}
		}
	} else {
		pwdPath := filepath.Join(keysDir, u.Hash()+".credential")
		if opts.Password != "" {
			data := []byte(u.Url(opts.Password))
			if err := fs.WriteFile(pwdPath, data, 0755); err != nil {
				return err
			}
		} else if hash != u.Hash() {
			oldPath := filepath.Join(keysDir, hash+".credential")
			if err := util.Copy(oldPath, pwdPath); err != nil {
				return err
			}
		}
	}
	for i := 0; i < len(s.c.Users); i++ {
		if s.c.Users[i].Hash() == hash {
			s.c.Users[i] = u
		}
	}
	return saveConfig(s.c)
}

func (s *Store) Delete(u *User) error {
	p, err := u.KeyPath()
	if err != nil {
		return err
	}
	if u.IsSSH() {
		os.Remove(p)
		os.Remove(p + ".pub")
	} else {
		os.Remove(p + ".credential")
	}
	users := make([]*User, 0)
	for _, v := range s.c.Users {
		if v.Hash() != u.Hash() {
			users = append(users, v)
		}
	}
	s.c.Users = users
	return saveConfig(s.c)
}

func (s *Store) List(filter string, strict bool) []*User {
	users := make([]*User, 0)
	for _, u := range s.c.Users {
		if !strict && strings.Contains(u.Name, filter) {
			users = append(users, u)
		}
		if strict && u.Name == filter {
			users = append(users, u)
		}
	}
	return users
}

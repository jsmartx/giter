package ssh

import (
	"github.com/jsmartx/giter/util"
	"github.com/kevinburke/ssh_config"
	fs "io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	cfg *ssh_config.Config
}

func loadConfig() *ssh_config.Config {
	p, err := util.JoinPath("~/.ssh/", "config")
	if err != nil {
		return nil
	}
	f, err := os.Open(p)
	defer f.Close()
	if err != nil {
		return nil
	}
	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil
	}
	return cfg
}

func saveConfig(cfg *ssh_config.Config) error {
	root, err := util.Mkdir("~/.ssh/")
	if err != nil {
		return err
	}
	b, err := cfg.MarshalText()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(root, "config")
	// write marshaled data to the file
	return fs.WriteFile(cfgPath, b, 0755)
}

func New() *Config {
	cfg := loadConfig()
	if cfg != nil {
		return &Config{cfg: cfg}
	}
	cfg = &ssh_config.Config{
		Hosts: make([]*ssh_config.Host, 0),
	}
	saveConfig(cfg)
	return &Config{cfg: cfg}
}

type Host struct {
	Key          string
	Hostname     string
	Port         string
	IdentityFile string
}

func (h *Host) Transform() (*ssh_config.Host, error) {
	pattern, err := ssh_config.NewPattern(h.Key)
	if err != nil {
		return nil, err
	}
	nodes := make([]ssh_config.Node, 0)
	if h.Hostname != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  HostName", Value: h.Hostname})
	}
	if h.Port != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  Port", Value: h.Port})
	}
	if h.IdentityFile != "" {
		nodes = append(nodes, &ssh_config.KV{Key: "  IdentityFile", Value: h.IdentityFile})
	}
	nodes = append(nodes, &ssh_config.Empty{})
	return &ssh_config.Host{
		Patterns:   []*ssh_config.Pattern{pattern},
		Nodes:      nodes,
		EOLComment: " -- added by giter",
	}, nil
}

func (c *Config) SetHost(h *Host) error {
	host, err := h.Transform()
	if err != nil {
		return err
	}
	for i, v := range c.cfg.Hosts {
		for _, pattern := range v.Patterns {
			if pattern.String() == h.Key {
				c.cfg.Hosts[i] = host
				return saveConfig(c.cfg)
			}
		}
	}
	c.cfg.Hosts = append(c.cfg.Hosts, host)
	return saveConfig(c.cfg)
}

package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io"
	fs "io/ioutil"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

func Keygen(publicPath, privatePath string, bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// generate and write private key as PEM
	privateFile, err := os.OpenFile(privatePath, os.O_RDWR|os.O_CREATE, 0600)
	defer privateFile.Close()
	if err != nil {
		return err
	}
	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, privatePEM); err != nil {
		return err
	}
	// generate and write public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return fs.WriteFile(publicPath, ssh.MarshalAuthorizedKey(publicKey), 0600)
}

func JoinHostPort(host, port string) string {
	if port != "" {
		return fmt.Sprintf("%s:%s", host, port)
	} else {
		return host
	}
}

func SplitHostPort(host string) (string, string) {
	h, p, err := net.SplitHostPort(host)
	if err != nil {
		return host, ""
	} else {
		return h, p
	}
}

var ScpRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)
var Schemes = []string{"git", "https", "http", "git+ssh", "ssh"}

func ParseURL(repo string) (*url.URL, error) {
	var err error
	var repoURL *url.URL
	if m := ScpRe.FindStringSubmatch(repo); m != nil {
		repoURL = &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   m[3],
		}
	} else {
		repoURL, err = url.Parse(repo)
		if err != nil {
			return nil, err
		}
	}
	for _, scheme := range Schemes {
		if repoURL.Scheme == scheme {
			return repoURL, nil
		}
	}
	return nil, errors.New(repoURL.Scheme + " is not supported")
}

type PromptConfig struct {
	Prompt  string
	Default string
	Silent  bool
}

func (c *PromptConfig) String() string {
	if c.Default != "" {
		return fmt.Sprintf("%s(%s) ", c.Prompt, c.Default)
	} else {
		return c.Prompt
	}
}

func CheckError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Error: %s", err))
	os.Exit(1)
}

func Prompt(cfg *PromptConfig) string {
	r, err := readline.NewEx(&readline.Config{
		Prompt:     cfg.String(),
		EnableMask: cfg.Silent,
		MaskRune:   42,
	})
	CheckError(err)
	defer r.Close()
	for {
		line, err := r.Readline()
		CheckError(err)
		if line != "" {
			return line
		}
		if cfg.Default != "" {
			return cfg.Default
		}
	}
}

func IsExist(path string) bool {
	p, err := homedir.Expand(path)
	if err != nil {
		return false
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func Mkdir(paths ...string) (string, error) {
	p, err := homedir.Expand(filepath.Join(paths...))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(p, 0755); err != nil {
		return "", err
	}
	return p, nil
}

func JoinPath(paths ...string) (string, error) {
	p, err := homedir.Expand(filepath.Join(paths...))
	if err != nil {
		return "", err
	}
	return p, nil
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func SysSSHConfig() (string, error) {
	keypairs := []string{"~/.ssh/id_rsa", "~/.ssh/id_ed25519"}
	for _, k := range keypairs {
		if IsExist(k) && IsExist(k+".pub") {
			return homedir.Expand(k)
		}
	}
	return "", errors.New("System ssh config not found")
}

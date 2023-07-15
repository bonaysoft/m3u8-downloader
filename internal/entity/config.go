package entity

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Profiles []Profile `json:"profiles"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	return cfg, viper.Unmarshal(cfg)
}

func (c *Config) SelectDecrypter(name string) *Decrypter {
	v, ok := lo.Find(c.Profiles, func(item Profile) bool { return item.Name == name })
	if !ok {
		return &Decrypter{Type: "none", Properties: Properties{}}
	}

	return &v.Decrypter
}

func (c *Config) DecrypterInstall(name string, decrypter *Decrypter) error {
	_, idx, ok := lo.FindIndexOf(c.Profiles, func(item Profile) bool { return item.Name == name })
	if ok {
		c.Profiles[idx].Name = name
		c.Profiles[idx].Decrypter = *decrypter
		fmt.Println("Reinstall the existing decrypter: ", name)
		return nil
	}

	c.Profiles = append(c.Profiles, Profile{Name: name, Decrypter: *decrypter})
	return nil
}

func (c *Config) DecrypterUninstall(name string) error {
	_, idx, ok := lo.FindIndexOf(c.Profiles, func(item Profile) bool { return item.Name == name })
	if !ok {
		return fmt.Errorf("not found the profile %s", name)
	}

	c.Profiles = append(c.Profiles[:idx], c.Profiles[idx+1:]...)
	return nil
}

func (c *Config) Save() error {
	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(c)
}

type Profile struct {
	Name      string    `json:"name"`
	Decrypter Decrypter `json:"decrypter"`
}

type Decrypter struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

func NewDecrypter(Type string, properties Properties) *Decrypter {
	return &Decrypter{Type: Type, Properties: properties}
}

type Properties struct {
	Cmd []string `json:"cmd,omitempty" yaml:"cmd,omitempty"`

	BaseURL string `json:"baseURL,omitempty" yaml:"baseURL,omitempty"`
	Token   string `json:"token,omitempty" yaml:"token,omitempty"`
}

func (p *Properties) String() string {
	if len(p.Cmd) > 0 {
		return fmt.Sprintf("cmd=[ %s ]", strings.Join(p.Cmd, " "))
	}

	return fmt.Sprintf("baseURL=%s, token=%s", p.BaseURL, p.Token)
}

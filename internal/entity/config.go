package entity

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
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

type Profile struct {
	Name      string    `json:"name"`
	Decrypter Decrypter `json:"decrypter"`
}

type Decrypter struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Cmd []string `json:"cmd,omitempty"`

	BaseURL string `json:"baseURL,omitempty"`
	Token   string `json:"token,omitempty"`
}

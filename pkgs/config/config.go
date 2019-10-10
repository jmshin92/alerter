package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigName = "config.yml"
	DefaultConfigDir  = "conf"
)

type Config struct {
	TargetUri string
	Alert     *Alert
	Mail      *Mail
}

var (
	c *Config
	m sync.Mutex
)

func NewConfig() *Config {
	conf := &Config{
		Alert: NewAlert(),
		Mail: NewMail(),
	}
	return conf
}

func GetConfig(args ...string) (*Config, error) {
	m.Lock()
	defer m.Unlock()

	if len(args) != 0 {
		conf, err := load(args[0])
		if err != nil {
			return nil, err
		}
		c = conf
	}
	return c, nil
}

func load(p string) (*Config, error) {
	conf := NewConfig()

	b, err := ioutil.ReadFile(p)
	if err != nil {
		err = errors.Wrapf(err, "failed to read config file[%v]", p)
		return nil, err
	}

	if err := yaml.Unmarshal(b, conf); err != nil {
		err = errors.Wrapf(err, "failed to unmarshal config file[%v]", p)
		return nil, err
	}

	return conf, nil
}

func Write(p string) error {
	m.Lock()
	defer m.Unlock()

	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if err := c.Write(p); err != nil {
		return err
	}
	return nil
}

func (this *Config) Write(p string) error {
	b, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(p, b, 0655); err != nil {
		return err
	}

	return nil
}
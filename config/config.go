package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

//Config - contains some configuration data
type Config struct {
	Page  string   `json:"page"`
	Name  string   `json:"policy_name"`
	Roles []string `json:"roles"`
}

var config *Config
var once sync.Once
var e error

//Init - initialize
func Init(cErr chan error) {
	once.Do(func() {
		config = &Config{}
		e = loadConfig()
	})
	cErr <- e
}

//Get - return copy of config
func Get() Config {
	return *config
}

func loadConfig() error {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("cannot find config file: %s\nPlease, add config file and restart program", err)
	}
	config = &Config{}
	if err = json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("corrupted data in config file: %s\nPlease, correct config and restart program", err)
	}
	return nil
}

//ContainsRole - return true if role is in the list, false otherwise
func (c Config) ContainsRole(role string) bool {
	for _, r := range c.Roles {
		if strings.ToLower(r) == strings.ToLower(role) {
			return true
		}
	}
	return false
}

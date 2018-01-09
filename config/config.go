package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type Config struct {
	Page  string
	Name  string
	Roles []string
}

var config *Config
var once sync.Once

func Get() Config {
	once.Do(loadConfig)
	return *config
}

func loadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Cannot find config file: %s", err)
	}
	config = &Config{
		Page: "Page",
		Name: "Policy technical Action Name",
	}
	if err = json.Unmarshal(data, &config.Roles); err != nil {
		log.Fatalf("Corrupted data in config file: %s", err)
	}
}

func (c Config) ContainsRole(role string) bool {
	for _, r := range c.Roles {
		if strings.ToLower(r) == strings.ToLower(role) {
			return true
		}
	}
	return false
}

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

//Config - contains some configuration data
type Config struct {
	Page          string `json:"page"`
	Name          string `json:"policy_name"`
	RolesBegin    string `json:"roles_begin"`
	RolesEnd      string `json:"roles_end"`
	Type          string `json:"type"`
	TechGroupName string `json:"technical_group_name"`
	DisplayName   string `json:"display_name"`
}

var config *Config
var once sync.Once
var e error

//Init - read config from file and send error in given channel
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

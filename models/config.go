package models

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Data      int64    `yaml:"data"`
	IPAddress []string `yaml:"ip_addresses"`
	Secret    string   `yaml:"secret"`
	SelfIp    string   `yaml:"self_ip"`
}

func GetConfig() (*Config, error) {
	t := Config{}
	buffer, err := ioutil.ReadFile("./config.yaml")
	err = yaml.Unmarshal(buffer, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func SaveConfig(config *Config) error {
	d, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if err := ioutil.WriteFile("./config.yaml", d, 0664); err != nil {
		return err
	}
	return nil
}

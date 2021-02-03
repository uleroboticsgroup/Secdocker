package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var ConfigFile string = "config.yml"

type GeneralConf struct {
	SecurityOptions       []string `yaml:"secopts,omitempty"`
	DropLinuxCapabilities []string `yaml:"capdrop,omitempty"`
	AddLinuxCapabilities  []string `yaml:"capadd,omitempty"`
	Memory                string   `yaml:"memory,omitempty"`
	CPU                   string   `yaml:"cpu,omitempty"`
	Environment           []string `yaml:"environment,omitempty"`
	User                  string   `yaml:"user,omitempty"`
}

type RestrictionsConf struct {
	Ports            []string `yaml:"ports"`
	Mounts           []string `yaml:"mounts"`
	Users            []string `yaml:"users"`
	Environment      []string `yaml:"environment"`
	SecurityPolicies []string `yaml:"securitypolicies"`
	Images           []string `yaml:"images"`
	Privileged       bool     `yaml:"privileged"`
}

// Config is the struct that holds all the restriction and general configurations
type Config struct {
	Plugins      []string         `yaml:"plugins"`
	DockerAPI    string           `yaml:"dockerapi"`
	General      GeneralConf      `yaml:"general"`
	Restrictions RestrictionsConf `yaml:"restrictions"`
}

// LoadConfig loads all the config from a yaml file and returns a Config object
func LoadConfig() Config {
	f, err := os.Open(ConfigFile)
	checkErr(err)
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	checkErr(err)

	return cfg
}

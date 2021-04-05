package config_test

import (
	"reflect"
	"testing"

	"github.com/uleroboticsgroup/Secdocker/config"
)

func TestLoadConfig(t *testing.T) {
	testConf := config.Config{
		DockerAPI: "http://localhost",
		General: config.GeneralConf{
			SecurityOptions:       []string{"opt1", "opt2"},
			DropLinuxCapabilities: []string{"cap1", "cap2"},
			AddLinuxCapabilities:  []string{"cap3", "cap4"},
			Memory:                "1g",
			CPU:                   "0.25",
			Environment:           []string{"MY_ENV=true", "MY_ENV2=1"},
			User:                  "1000",
		},
		Restrictions: config.RestrictionsConf{
			Ports:            []string{"8080", "3000"},
			Mounts:           []string{"/root"},
			Users:            []string{"root"},
			Environment:      []string{"USER=0"},
			SecurityPolicies: []string{"privileged"},
			Privileged:       true,
		},
	}

	config.ConfigFile = "../testing/testconfig.yml"
	conf := config.LoadConfig()
	if !reflect.DeepEqual(conf, testConf) {
		t.Fail()
	}
}

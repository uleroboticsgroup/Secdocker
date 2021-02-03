package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type plugin string

type config struct {
	ScriptPath string `yaml:"scriptpath"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (p plugin) Process(image string) bool {
	log.Info("Notary image analysis initialized")

	f, err := os.Open("./plugins/notary/config.yml")
	checkErr(err)
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	checkErr(err)

	args := strings.Split(image, ":")
	var name, tag string
	if len(args) == 1 {
		name = args[0]
		tag = "latest"
	} else if len(args) == 2 {
		name = args[0]
		tag = args[1]
	} else {
		fmt.Print("Image invalid")
		return false
	}

	cmd := exec.Command("/bin/bash", cfg.ScriptPath, name, tag)
	err = cmd.Run()
	err = cmd.Wait()
	returnCode := cmd.ProcessState.ExitCode()

	switch returnCode {
	case 0:
		log.Info("They match")
		return true
	case 1:
		log.Error("They don't match")
		return false
	case 2:
		log.Info("Image is not present locally")
		return true
	default:
		log.Error("Unexpected return code" + fmt.Sprint(returnCode))
		output, _ := cmd.Output()
		log.Error(output)
		return true
	}
}

var Plugin plugin

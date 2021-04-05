package docker

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/uleroboticsgroup/Secdocker/config"
	"github.com/uleroboticsgroup/Secdocker/plugins"
)

type ContainerOpts struct {
	Mounts           []string
	Ports            []string
	Env              []string
	SecurityPolicies []string
	Image            string
	Entrypoint       string
	User             string
	Privileged       bool
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func solveCollisions(first, second []string, separator string) []string {
	unique := map[string]string{}
	finalList := []string{}

	for _, v := range second {
		if separator != "" {
			unique[strings.Split(v, separator)[0]] = strings.Split(v, separator)[1]
		} else {
			unique[v] = ""
		}
	}

	for _, v := range first {
		if separator != "" {
			unique[strings.Split(v, separator)[0]] = strings.Split(v, separator)[1]
		} else {
			unique[v] = ""
		}
	}

	for k, v := range unique {
		finalList = append(finalList, k+separator+v)
	}

	return finalList
}

// AddGeneralRestrictions adds the restrictions from config to the already container creation instructions
func AddGeneralRestrictions(data ContainerOpts) ContainerOpts {
	configData := config.LoadConfig()

	if configData.General.User != "" {
		data.User = configData.General.User
	}
	data.Env = solveCollisions(configData.General.Environment, data.Env, "=")

	return data
}

// ProcessAPICreateRequest checks if the data pass all the restrictions
func ProcessAPICreateRequest(data ContainerOpts) bool {
	if plugins.ProcessPlugins(data.Image) && CheckPermissions(data) {
		return true
	}

	return false
}

// ExecuteCommand executes a command locally (used by the CLI version of the program)
func ExecuteCommand(command string, args []string) string {
	log.Info("Running [", command, " ", strings.Join(args, " "), "]")
	cmd := exec.Command(command, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(out))
	return string(out)
}

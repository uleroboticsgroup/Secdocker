package commandline

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"niebla.unileon.es/DavidFerng/secdocker/config"
	"niebla.unileon.es/DavidFerng/secdocker/docker-utilities"
	dockerutilities "niebla.unileon.es/DavidFerng/secdocker/docker-utilities"
)

type Flag string

const (
	Port       Flag = "-p"
	Mount      Flag = "-v"
	Env        Flag = "-e"
	Entrypoint Flag = "--entrypoint"
	User       Flag = "-u"
	None       Flag = ""
)

func cleanEqualArgs(args []string) []string {
	finalArgs := []string{}
	isFlag := false

	for _, arg := range args {
		if !isFlag && strings.Contains(arg, "=") {
			finalArgs = append(finalArgs, strings.Split(arg, "=")[0])
			finalArgs = append(finalArgs, strings.Split(arg, "=")[1])
		} else {
			finalArgs = append(finalArgs, arg)
			isFlag = !isFlag
		}
	}
	return finalArgs
}

func GenerateArgsFromConfig() []string {
	config := config.LoadConfig()
	args := []string{}

	if config.General.CPU != "" {
		args = append(args, "--cpus")
		args = append(args, config.General.CPU)
	}

	if config.General.User != "" {
		args = append(args, "-u")
		args = append(args, config.General.User)
	}

	if config.General.Memory != "" {
		args = append(args, "-m")
		args = append(args, config.General.Memory)
	}

	for _, item := range config.General.Environment {
		args = append(args, "-e")
		args = append(args, item)
	}

	for _, item := range config.General.AddLinuxCapabilities {
		args = append(args, "--cap-add")
		args = append(args, item)
	}

	for _, item := range config.General.DropLinuxCapabilities {
		args = append(args, "--cap-drop")
		args = append(args, item)
	}

	for _, item := range config.General.SecurityOptions {
		args = append(args, "--security-opt")
		args = append(args, item)
	}

	return args
}

func ParseRunArgs(args []string) dockerutilities.ContainerOpts {
	var containerOpts dockerutilities.ContainerOpts
	activeFlag := None

	args = cleanEqualArgs(args)
	for _, arg := range args {
		switch arg {
		case "-p":
			activeFlag = Port
		case "-P":
			activeFlag = Port
		case "-v":
			activeFlag = Mount
		case "--volume":
			activeFlag = Mount
		case "-e":
			activeFlag = Env
		case "-u":
			activeFlag = User
		case "--user":
			activeFlag = User
		case "--entrypoint":
			activeFlag = Entrypoint
		default:
			switch activeFlag {
			case Port:
				// format: ip:hostPort:containerPort | ip::containerPort | hostPort:containerPort | containerPort
				ports := strings.Split(arg, ":")
				if len(ports) == 3 {
					containerOpts.Ports = append(containerOpts.Ports, ports[1])
				} else if len(ports) == 2 {
					containerOpts.Ports = append(containerOpts.Ports, ports[0])
				}
			case Mount:
				mount := strings.Split(arg, ":")
				containerOpts.Mounts = append(containerOpts.Mounts, mount[0])
			case Env:
				containerOpts.Env = append(containerOpts.Env, arg)
			case User:
				containerOpts.User = arg
			case Entrypoint:
				containerOpts.Entrypoint = arg
			case None:
				containerOpts.Image = arg
			}
			activeFlag = None
		}
	}

	return containerOpts
}

func ParseCommandLineArgs(args []string) dockerutilities.ContainerOpts {
	if args[0] == "run" {
		return ParseRunArgs(args[1:])
	}

	return dockerutilities.ContainerOpts{}
}

func StartCLI() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: " + os.Args[0] + " run (docker options) image")
		os.Exit(1)
	}

	log.Info("Received ", os.Args[1:])
	arguments := ParseCommandLineArgs(os.Args[1:])
	if docker.CheckPermissions(arguments) {
		log.Info("Executing...")
		addedArgs := GenerateArgsFromConfig()
		addedArgs = append(addedArgs, "-d")

		finalArgs := append([]string{os.Args[0]}, append(addedArgs, os.Args[1:]...)...)
		docker.ExecuteCommand("docker", finalArgs)
	} else {
		log.Error("Command contains invalid or forbidden values. Aborting.")
	}
}

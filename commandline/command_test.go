package commandline_test

import (
	"reflect"
	"testing"

	"niebla.unileon.es/DavidFerng/secdocker/commandline"
	"niebla.unileon.es/DavidFerng/secdocker/config"
	"niebla.unileon.es/DavidFerng/secdocker/docker-utilities"
)

func indexesOfStringInSlice(target string, list []string) []int {
	indexes := []int{}
	for i, item := range list {
		if item == target {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func TestParseCommandLineArgs(t *testing.T) {
	testCases := []struct {
		args   []string
		result docker.ContainerOpts
	}{
		{[]string{"run"}, docker.ContainerOpts{}},

		{[]string{"run", "-p", "8080:8081"}, docker.ContainerOpts{Ports: []string{"8080"}}},

		{[]string{"error", "error", "error"}, docker.ContainerOpts{}},
	}

	for _, testCase := range testCases {
		result := commandline.ParseCommandLineArgs(testCase.args)
		if !reflect.DeepEqual(result, testCase.result) {
			t.Errorf("Expected %+v, found %+v", testCase.result, result)
		}
	}
}

func TestParseRunArgs(t *testing.T) {
	testCases := []struct {
		args   []string
		result docker.ContainerOpts
	}{
		{[]string{}, docker.ContainerOpts{}},

		{[]string{"-p", "8080:8081"}, docker.ContainerOpts{Ports: []string{"8080"}}},
		{[]string{"-p", "8080:8081", "-p", "3000:3001"}, docker.ContainerOpts{Ports: []string{"8080", "3000"}}},
		{[]string{"-p=8080:8081"}, docker.ContainerOpts{Ports: []string{"8080"}}},

		{[]string{"-v", "/root:/tmp"}, docker.ContainerOpts{Mounts: []string{"/root"}}},
		{[]string{"-v", "/root:/tmp", "-v", "/usr:/share"}, docker.ContainerOpts{Mounts: []string{"/root", "/usr"}}},
		{[]string{"-v=/root:/tmp", "-v", "/usr:/share"}, docker.ContainerOpts{Mounts: []string{"/root", "/usr"}}},

		{[]string{"-e", "TEST=1"}, docker.ContainerOpts{Env: []string{"TEST=1"}}},
		{[]string{"-e", "TEST=1", "-e", "TEST=2"}, docker.ContainerOpts{Env: []string{"TEST=1", "TEST=2"}}},

		{[]string{"test"}, docker.ContainerOpts{Image: "test"}},

		{[]string{"-u", "root"}, docker.ContainerOpts{User: "root"}},

		{[]string{"--entrypoint", "/bin/bash"}, docker.ContainerOpts{Entrypoint: "/bin/bash"}},
	}

	for _, testCase := range testCases {
		result := commandline.ParseRunArgs(testCase.args)
		if !reflect.DeepEqual(result, testCase.result) {
			t.Errorf("Expected %+v, found %+v", testCase.result, result)
		}
	}
}

func TestGenerateArgsFromConfig(t *testing.T) {
	testConf := []string{"--security-opt", "opt1", "--security-opt", "opt2",
		"--cap-drop", "cap1", "--cap-drop", "cap2",
		"--cap-add", "cap3", "--cap-add", "cap4",
		"-m", "1g",
		"--cpus", "0.25", "-e", "MY_ENV=true", "-e", "MY_ENV2=1",
		"-u", "1000"}

	config.ConfigFile = "../testing/testconfig.yml"
	args := commandline.GenerateArgsFromConfig()

	for current := 0; current < len(testConf); current += 2 {
		indexes := indexesOfStringInSlice(testConf[current], args)
		found := false
		for _, index := range indexes {
			if testConf[current+1] == args[index+1] {
				found = true
			}
		}

		if !found {
			t.Errorf("Tag %s: %s not found", testConf[current], testConf[current+1])
		}
	}
}

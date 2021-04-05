package docker_test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/uleroboticsgroup/Secdocker/config"
	"github.com/uleroboticsgroup/Secdocker/docker-utilities"
)

func TestCheckIntersection(t *testing.T) {
	log.SetOutput(os.Stdout)

	testCases := []struct {
		userItems     []string
		securityItems []string
		result        bool
	}{
		{[]string{"test1", "test2", "test3"}, []string{"test4", "test1", "test5"}, true},
		{[]string{"test1", "test2", "test3"}, []string{"test4", "test5", "test6"}, false},
	}

	for _, testCase := range testCases {
		result := docker.CheckIntersection(testCase.userItems, testCase.securityItems, "")
		if result != testCase.result {
			t.Errorf("Expected %t, found %t", testCase.result, result)
		}
	}
}

func TestCheckPermissions(t *testing.T) {
	config.ConfigFile = "../testing/testconfig.yml"
	log.SetOutput(os.Stdout)

	testCases := []struct {
		parameters docker.ContainerOpts
		result     bool
	}{
		{
			parameters: docker.ContainerOpts{
				Ports:      []string{"8080", "3000"},
				Mounts:     []string{"/usr"},
				Env:        []string{"USER=1", "USER=2"},
				Privileged: false,
			}, result: false,
		},
		{

			parameters: docker.ContainerOpts{
				Ports:      []string{"8081", "3001"},
				Mounts:     []string{"/root"},
				Env:        []string{"USER=1", "USER=2"},
				Privileged: false,
			}, result: false,
		},
		{
			parameters: docker.ContainerOpts{
				Ports:      []string{"8081", "3001"},
				Mounts:     []string{"/usr"},
				Env:        []string{"USER=0", "USER=1"},
				Privileged: false,
			},
			result: false,
		},
		{
			parameters: docker.ContainerOpts{
				Ports:      []string{"8081", "3001"},
				Mounts:     []string{"/usr"},
				Env:        []string{"USER=1", "USER=2"},
				Privileged: true,
			}, result: false,
		},
		{
			parameters: docker.ContainerOpts{
				Ports:      []string{"8081", "3001"},
				Mounts:     []string{"/usr"},
				Env:        []string{"USER=1", "USER=2"},
				Privileged: false,
			}, result: true,
		},
	}

	for _, testCase := range testCases {
		result := docker.CheckPermissions(testCase.parameters)
		if result != testCase.result {
			t.Errorf("Expected %t, found %t", testCase.result, result)
		}
	}
}

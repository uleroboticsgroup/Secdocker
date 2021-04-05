package docker_test

import (
	"testing"

	. "github.com/uleroboticsgroup/Secdocker/docker-utilities"

	"github.com/uleroboticsgroup/Secdocker/config"
)

func TestExecuteCommand(t *testing.T) {
	config.ConfigFile = "../testing/testconfig.yml"
	out := ExecuteCommand("echo", []string{"test test"})
	outputString := "test test\n"
	if out != outputString {
		t.Errorf("Expected %s, found %s", outputString, out)
	}
}

func TestSolveCollisions(t *testing.T) {
	testCases := []struct {
		first     []string
		second    []string
		result    []string
		separator string
	}{
		{
			second:    []string{"a=a", "b=b", "c=c"},
			first:     []string{"d=d", "e=e", "f=f"},
			result:    []string{"a=a", "b=b", "c=c", "d=d", "e=e", "f=f"},
			separator: "=",
		},
		{
			second:    []string{"a=a", "b=b", "c=c"},
			first:     []string{"a=d", "b=e", "c=f"},
			result:    []string{"a=d", "b=e", "c=f"},
			separator: "=",
		},
		{
			second:    []string{"a=a", "b=b", "c=c"},
			first:     []string{"a=d", "b=e", "c=f"},
			result:    []string{"a=a", "b=b", "c=c", "a=d", "b=e", "c=f"},
			separator: "",
		},
	}

	for _, item := range testCases {
		testResult := SolveCollisions(item.first, item.second, item.separator)
		areEqual := true

		if len(item.result) != len(testResult) {
			areEqual = false
		}

		for _, v := range item.result {
			found := false
			for _, v2 := range testResult {
				if v == v2 {
					found = true
				}
			}

			if !found {
				areEqual = false
			}
		}

		if !areEqual {
			t.Errorf("Expected %s, found %s", item.result, testResult)
		}

	}
}

func TestAddGeneralRestrictions(t *testing.T) {
	config.ConfigFile = "../testing/testconfig.yml"
	data := ContainerOpts{User: "root", Env: []string{"MY_ENV=false"}}
	exptectedData := ContainerOpts{User: "1000", Env: []string{"MY_ENV=true", "MY_ENV2=1"}}
	data = AddGeneralRestrictions(data)

	if data.User != exptectedData.User {
		t.Errorf("Expected %s, found %s", exptectedData.User, data.User)
	}

	for _, i := range exptectedData.Env {
		found := false
		for _, i2 := range data.Env {
			if i == i2 {
				found = true
			}
		}

		if !found {
			t.Errorf("%s not found", i)
		}
	}
}

func TestProcessAPiCreateRequest(t *testing.T) {
	config.ConfigFile = "../testing/testconfig.yml"

	testCases := []struct {
		data   ContainerOpts
		result bool
	}{
		{
			data:   ContainerOpts{User: "root"},
			result: false,
		},
		{
			data:   ContainerOpts{User: "1000"},
			result: true,
		},
	}

	for _, testCase := range testCases {
		if result := ProcessAPICreateRequest(testCase.data); testCase.result != result {
			t.Errorf("Expected %t; found %t; %+v", testCase.result, result, testCase.data)
		}
	}

}

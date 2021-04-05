package httpserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/uleroboticsgroup/Secdocker/config"
	. "github.com/uleroboticsgroup/Secdocker/httpserver"
)

/*
var postToApi = apitest.NewMock().
	Post("http://localhost/containers/create").
	RespondWith().
	Body(`{"message: "ok"}`).
	Status(http.StatusOK).
	End()

var postToApiImages = apitest.NewMock().
	Get("http://localhost/images/list").
	RespondWith().
	Body(`{"message: "ok"}`).
	Status(http.StatusOK).
	End()
*/

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	fmt.Println(string(b))
	fmt.Println("")
	fmt.Println(string(a))
	for i, v := range a {
		if v != b[i] {
			//du := b[i]
			return false
		}
	}
	return true
}

func TestProcessCreateContainer(t *testing.T) {
	config.ConfigFile = "../testing/testconfig.yml"

	jsonFile, err := os.Open("../testing/createcontainer.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	data, _ := ioutil.ReadAll(jsonFile)
	dataResult := []byte(`{"Hostname":"","Domainname":"","User":"1000","AttachStdin":false,"AttachStdout":true,"AttachStderr":true,"Tty":false,"OpenStdin":false,"StdinOnce":false,"Env":["FOO=bar","BAZ=quux","MY_ENV=true","MY_ENV2=1"],"Cmd":["date"],"Image":"ubuntu:18.04","Labels":{"com.example.license":"GPL","com.example.vendor":"Acme","com.example.version":"1.0"},"Volumes":{"/volumes/data":{}},"WorkingDir":"","NetworkDisabled":false,"MacAddress":"12:34:56:78:9a:bc","ExposedPorts":{"22/tcp":{}},"StopSignal":"SIGTERM","StopTimeout":10,"HostConfig":{"Binds":["/tmp:/tmp"],"Links":[],"Memory":0,"MemorySwap":0,"MemoryReservation":0,"KernelMemory":0,"NanoCPUs":0,"CpuPercent":0,"CpuShares":0,"CpuPeriod":0,"CpuRealtimePeriod":0,"CpuRealtimeRuntime":0,"CpuQuota":0,"CpusetCpus":"","CpusetMems":"","MaximumIOps":0,"MaximumIOBps":0,"BlkioWeight":300,"BlkioWeightDevice":[{"Path":"","Weight":0}],"BlkioDeviceReadBps":[{"Path":"","Rate":0}],"BlkioDeviceReadIOps":[{"Path":"","Rate":0}],"BlkioDeviceWriteBps":[{"Path":"","Rate":0}],"BlkioDeviceWriteIOps":[{"Path":"","Rate":0}],"DeviceRequests":[],"MemorySwappiness":60,"OomKillDisable":false,"OomScoreAdj":500,"PidMode":"","PidsLimit":0,"PortBindings":{"22/tcp":[{"HostIp":"0.0.0.0","HostPort":"11022"}],"80/tcp":[{"HostIp":"","HostPort":"8081"}]},"PublishAllPorts":false,"Privileged":false,"ReadonlyRootfs":false,"Dns":["8.8.8.8"],"DnsOptions":[""],"DnsSearch":[""],"VolumesFrom":[],"CapAdd":["NET_ADMIN"],"CapDrop":["MKNOD"],"GroupAdd":["newgroup"],"RestartPolicy":{"Name":"","MaximumRetryCount":0},"AutoRemove":true,"NetworkMode":"bridge","Devices":[],"Ulimits":[{"Name":"","Soft":0,"Hard":0}],"LogConfig":{"Type":"json-file","Config":{}},"SecurityOpt":[],"StorageOpt":null,"CgroupParent":"","VolumeDriver":"","ShmSize":67108864},"NetworkingConfig":{"EndpointsConfig":{"isolated_nw":{"IPAMConfig":{"IPv4Address":"172.20.30.33","IPv6Address":"2001:db8:abcd::3033","LinkLocalIPs":["169.254.34.68","fe80::3468"]},"Links":["container_1","container_2"],"Aliases":["server_x","server_y"],"NetworkID":"","EndpointID":"","Gateway":"","IPAdress":"","IPPrefixLen":0,"IPv6Gateway":"","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"MadAddress":"","DriverOpts":null}}}}`)

	jsonFile, err = os.Open("../testing/createcontainerforbidden.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	forbiddenData, _ := ioutil.ReadAll(jsonFile)

	type TestStruct struct {
		runOpts []byte
		result  []byte
	}

	tests := []TestStruct{
		TestStruct{runOpts: data, result: dataResult},
		TestStruct{runOpts: forbiddenData, result: []byte{}},
	}

	for _, test := range tests {
		req := http.Request{}

		stringReader := bytes.NewReader(test.runOpts)
		stringReadCloser := ioutil.NopCloser(stringReader)
		req.Body = stringReadCloser
		req.ContentLength = int64(len(test.runOpts))

		req.Method = "POST"
		req.Proto = "HTTP/1.1"
		req.Body = stringReadCloser

		if data := ProcessCreateContainer(&req); len(data) != len(test.result) {

			fmt.Println(string(data))
			fmt.Println(string(test.result))
			t.Errorf("Expected length: %d; got %d", len(test.result), len(data))
		}
	}
}

func TestCreateOptsFromApiData(t *testing.T) {
	jsonFile, err := os.Open("../testing/createcontainer.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValues, _ := ioutil.ReadAll(jsonFile)
	var rawOpts RunOpts
	err = json.Unmarshal(byteValues, &rawOpts)
	containerOpts := CreateOptsFromAPIData(rawOpts)

	for i, _ := range containerOpts.Mounts {
		if containerOpts.Mounts[i] != rawOpts.HostConfig.Binds[i] {
			t.Errorf("Mounts: Expected: %s; got %s", rawOpts.HostConfig.Binds[i], containerOpts.Mounts[i])
		}
	}

	for i, _ := range containerOpts.Env {
		if containerOpts.Env[i] != rawOpts.Env[i] {
			t.Errorf("Env: Expected: %s; got %s", rawOpts.Env[i], containerOpts.Env[i])
		}
	}

	for i, _ := range containerOpts.SecurityPolicies {
		if containerOpts.SecurityPolicies[i] != rawOpts.HostConfig.SecurityOpt[i] {
			t.Errorf("Env: Expected: %s; got %s", rawOpts.HostConfig.SecurityOpt[i], containerOpts.SecurityPolicies[i])
		}
	}

	if containerOpts.Entrypoint != rawOpts.Entrypoint {
		t.Errorf("Env: Expected: %s; got %s", rawOpts.Entrypoint, containerOpts.Entrypoint)
	}

	if containerOpts.Image != rawOpts.Image {
		t.Errorf("Env: Expected: %s; got %s", rawOpts.Image, containerOpts.Image)
	}

	if containerOpts.User != rawOpts.User {
		t.Errorf("Env: Expected: %s; got %s", rawOpts.User, containerOpts.User)
	}

	if containerOpts.Privileged != rawOpts.HostConfig.Privileged {
		t.Errorf("Env: Expected: %t; got %t", rawOpts.HostConfig.Privileged, containerOpts.Privileged)
	}

	ports := []string{"0.0.0.0:11022", "8081"}
	for _, v := range containerOpts.Ports {
		found := false
		for _, port := range ports {
			if port == v {
				found = true
			}
		}

		if !found {
			t.Errorf("Env: Expected: %v not found", v)
		}
	}

}

func TestCreateRunDataFromOpts(t *testing.T) {
	jsonFile, err := os.Open("../testing/createcontainer.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValues, _ := ioutil.ReadAll(jsonFile)
	var rawOpts RunOpts
	var rawOpts2 RunOpts
	err = json.Unmarshal(byteValues, &rawOpts)
	containerOpts := CreateOptsFromAPIData(rawOpts)
	rawOpts2 = CreateRunDataFromOpts(rawOpts2, containerOpts)

	for i, _ := range containerOpts.Mounts {
		if containerOpts.Mounts[i] != rawOpts2.HostConfig.Binds[i] {
			t.Errorf("Mounts: Expected: %s; got %s", containerOpts.Mounts[i], rawOpts2.HostConfig.Binds[i])
		}
	}

	for i, _ := range containerOpts.Env {
		if containerOpts.Env[i] != rawOpts2.Env[i] {
			t.Errorf("Env: Expected: %s; got %s", containerOpts.Env[i], rawOpts2.Env[i])
		}
	}

	for i, _ := range containerOpts.SecurityPolicies {
		if containerOpts.SecurityPolicies[i] != rawOpts2.HostConfig.SecurityOpt[i] {
			t.Errorf("Env: Expected: %s; got %s", containerOpts.SecurityPolicies[i], rawOpts2.HostConfig.SecurityOpt[i])
		}
	}

	if containerOpts.Entrypoint != rawOpts2.Entrypoint {
		t.Errorf("Env: Expected: %s; got %s", containerOpts.Entrypoint, rawOpts2.Entrypoint)
	}

	if containerOpts.Image != rawOpts2.Image {
		t.Errorf("Env: Expected: %s; got %s", containerOpts.Image, rawOpts2.Image)
	}

	if containerOpts.User != rawOpts2.User {
		t.Errorf("Env: Expected: %s; got %s", containerOpts.User, rawOpts2.User)
	}

	if containerOpts.Privileged != rawOpts2.HostConfig.Privileged {
		t.Errorf("Env: Expected: %t; got %t", containerOpts.Privileged, rawOpts2.HostConfig.Privileged)
	}

}

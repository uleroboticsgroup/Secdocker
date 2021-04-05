package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/uleroboticsgroup/Secdocker/docker-utilities"
)

var dockerEngineVersion = "v1.24"

// RunOpts is the specification of the package of the Docker API
type RunOpts struct {
	Hostname        string              `json:"Hostname"`
	Domainname      string              `json:"Domainname"`
	User            string              `json:"User"`
	AttachStdin     bool                `json:"AttachStdin"`
	AttachStdout    bool                `json:"AttachStdout"`
	AttachStderr    bool                `json:"AttachStderr"`
	Tty             bool                `json:"Tty"`
	OpenStdin       bool                `json:"OpenStdin"`
	StdinOnce       bool                `json:"StdinOnce"`
	Env             []string            `json:"Env"`
	Cmd             []string            `json:"Cmd"`
	Entrypoint      string              `json:"Entrypoint,omitempty"`
	Image           string              `json:"Image"`
	Labels          map[string]string   `json:"Labels"`
	Volumes         map[string]struct{} `json:"Volumes"`
	WorkingDir      string              `json:"WorkingDir"`
	NetworkDisabled bool                `json:"NetworkDisabled"`
	MacAddress      string              `json:"MacAddress"`
	ExposedPorts    map[string]struct{} `json:"ExposedPorts"`
	StopSignal      string              `json:"StopSignal"`
	StopTimeout     int                 `json:"StopTimeout"`
	HostConfig      struct {
		Binds              []string `json:"Binds"`
		Links              []string `json:"Links"`
		Memory             int      `json:"Memory"`
		MemorySwap         int      `json:"MemorySwap"`
		MemoryReservation  int      `json:"MemoryReservation"`
		KernelMemory       int      `json:"KernelMemory"`
		NanoCPUs           int      `json:"NanoCPUs"`
		CPUPercent         int      `json:"CpuPercent"`
		CPUShares          int      `json:"CpuShares"`
		CPUPeriod          int      `json:"CpuPeriod"`
		CPURealtimePeriod  int      `json:"CpuRealtimePeriod"`
		CPURealtimeRuntime int      `json:"CpuRealtimeRuntime"`
		CPUQuota           int      `json:"CpuQuota"`
		CpusetCpus         string   `json:"CpusetCpus"`
		CpusetMems         string   `json:"CpusetMems"`
		MaximumIOps        int      `json:"MaximumIOps"`
		MaximumIOBps       int      `json:"MaximumIOBps"`
		BlkioWeight        int      `json:"BlkioWeight"`
		BlkioWeightDevice  []struct {
			Path   string `json:"Path"`
			Weight int    `json:"Weight"`
		} `json:"BlkioWeightDevice"`
		BlkioDeviceReadBps []struct {
			Path string `json:"Path"`
			Rate int    `json:"Rate"`
		} `json:"BlkioDeviceReadBps"`
		BlkioDeviceReadIOps []struct {
			Path string `json:"Path"`
			Rate int    `json:"Rate"`
		} `json:"BlkioDeviceReadIOps"`
		BlkioDeviceWriteBps []struct {
			Path string `json:"Path"`
			Rate int    `json:"Rate"`
		} `json:"BlkioDeviceWriteBps"`
		BlkioDeviceWriteIOps []struct {
			Path string `json:"Path"`
			Rate int    `json:"Rate"`
		} `json:"BlkioDeviceWriteIOps"`
		DeviceRequests []struct {
			Driver       string            `json:"Driver"`
			Count        int               `json:"Count"`
			DeviceIDs    []string          `json:"DeviceIDs"`
			Capabilities [][]string        `json:"Capabilities"`
			Options      map[string]string `json:"Options"`
		} `json:"DeviceRequests"`
		MemorySwappiness int    `json:"MemorySwappiness"`
		OomKillDisable   bool   `json:"OomKillDisable"`
		OomScoreAdj      int    `json:"OomScoreAdj"`
		PidMode          string `json:"PidMode"`
		PidsLimit        int    `json:"PidsLimit"`
		PortBindings     map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"PortBindings"`
		PublishAllPorts bool     `json:"PublishAllPorts"`
		Privileged      bool     `json:"Privileged"`
		ReadonlyRootfs  bool     `json:"ReadonlyRootfs"`
		DNS             []string `json:"Dns"`
		DNSOptions      []string `json:"DnsOptions"`
		DNSSearch       []string `json:"DnsSearch"`
		VolumesFrom     []string `json:"VolumesFrom"`
		CapAdd          []string `json:"CapAdd"`
		CapDrop         []string `json:"CapDrop"`
		GroupAdd        []string `json:"GroupAdd"`
		RestartPolicy   struct {
			Name              string `json:"Name"`
			MaximumRetryCount int    `json:"MaximumRetryCount"`
		} `json:"RestartPolicy"`
		AutoRemove  bool          `json:"AutoRemove"`
		NetworkMode string        `json:"NetworkMode"`
		Devices     []interface{} `json:"Devices"`
		Ulimits     []struct {
			Name string `json:"Name"`
			Soft int    `json:"Soft"`
			Hard int    `json:"Hard"`
		} `json:"Ulimits"`
		LogConfig struct {
			Type   string `json:"Type"`
			Config struct {
			} `json:"Config"`
		} `json:"LogConfig"`
		SecurityOpt  []string          `json:"SecurityOpt"`
		StorageOpt   map[string]string `json:"StorageOpt"`
		CgroupParent string            `json:"CgroupParent"`
		VolumeDriver string            `json:"VolumeDriver"`
		ShmSize      int               `json:"ShmSize"`
	} `json:"HostConfig"`
	NetworkingConfig struct {
		EndpointsConfig map[string]struct {
			IPAMConfig struct {
				IPv4Address  string   `json:"IPv4Address"`
				IPv6Address  string   `json:"IPv6Address"`
				LinkLocalIPs []string `json:"LinkLocalIPs"`
			} `json:"IPAMConfig"`
			Links               []string          `json:"Links"`
			Aliases             []string          `json:"Aliases"`
			NetworkID           string            `json:"NetworkID"`
			EndpointID          string            `json:"EndpointID"`
			Gateway             string            `json:"Gateway"`
			IPAdress            string            `json:"IPAdress"`
			IPPrefixLen         int               `json:"IPPrefixLen"`
			IPv6Gateway         string            `json:"IPv6Gateway"`
			GlobalIPv6Address   string            `json:"GlobalIPv6Address"`
			GlobalIPv6PrefixLen int64             `json:"GlobalIPv6PrefixLen"`
			MadAddress          string            `json:"MadAddress"`
			DriverOpts          map[string]string `json:"DriverOpts"`
		} `json:"EndpointsConfig"`
	} `json:"NetworkingConfig"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createOptsFromAPIData(rawOpts RunOpts) docker.ContainerOpts {
	containerOpts := docker.ContainerOpts{}
	containerOpts.Mounts = rawOpts.HostConfig.Binds
	containerOpts.Env = rawOpts.Env
	containerOpts.Entrypoint = rawOpts.Entrypoint
	containerOpts.Image = rawOpts.Image
	containerOpts.SecurityPolicies = rawOpts.HostConfig.SecurityOpt
	containerOpts.User = rawOpts.User
	containerOpts.Privileged = rawOpts.HostConfig.Privileged

	for _, v := range rawOpts.HostConfig.PortBindings {
		for _, item := range v {
			if item.HostIP != "" {
				containerOpts.Ports = append(containerOpts.Ports, item.HostIP+":"+item.HostPort)
			} else {
				containerOpts.Ports = append(containerOpts.Ports, item.HostPort)
			}
		}
	}

	return containerOpts
}

func createRunDataFromOpts(rawOpts RunOpts, opts docker.ContainerOpts) RunOpts {
	rawOpts.HostConfig.Binds = opts.Mounts
	rawOpts.Env = opts.Env
	rawOpts.Entrypoint = opts.Entrypoint
	rawOpts.Image = opts.Image
	rawOpts.HostConfig.SecurityOpt = opts.SecurityPolicies
	rawOpts.User = opts.User
	rawOpts.HostConfig.Privileged = opts.Privileged

	if len(rawOpts.Cmd) == 0 {
		rawOpts.Cmd = nil
	}

	return rawOpts
}

// ProcessCreateContainer handles a new http request to create a new container
func ProcessCreateContainer(req *http.Request) []byte {
	log.Info("New create request received: " + req.Method + " " + req.RequestURI)
	var rawOpts RunOpts
	body, err := ioutil.ReadAll(req.Body)
	checkErr(err)
	if len(body) == 0 {
		rawOpts.Image = req.URL.Query()["fromImage"][0] + ":" + req.URL.Query()["tag"][0]
	} else {
		err = json.Unmarshal(body, &rawOpts)
		checkErr(err)
	}

	if opts := createOptsFromAPIData(rawOpts); docker.ProcessAPICreateRequest(opts) {
		log.Info("Request is valid")
		opts := docker.AddGeneralRestrictions(opts)
		finalData := createRunDataFromOpts(rawOpts, opts)
		rawDataToSend, err := json.Marshal(finalData)
		checkErr(err)
		return rawDataToSend

	} else {
		log.Info("Request is not valid")
	}

	return []byte{}
}

/*
#########################################
## THIS CODE ONLY ACCEPTS HTTP CONNECTION.
## Docker CLI uses special TCP packets to work properly
#########################################

func sendToAPI(method, url string, data []byte) *http.Response {
	//https://docs.docker.com/engine/api/v1.40/#
	dockerAPI := config.LoadConfig().DockerAPI
	client := &http.Client{Timeout: 3 * time.Second}

	if dockerAPI[0] == '/' {
		fd := func(proto, addr string) (conn net.Conn, err error) {
			return net.Dial("unix", dockerAPI)
		}
		tr := &http.Transport{
			Dial: fd,
		}
		client = &http.Client{Transport: tr, Timeout: 3 * time.Second}
		url = "http://docker" + url
	} else {
		url = dockerAPI + url
	}

	fmt.Println(method, url)
	var resp *http.Response
	if method == "POST" || method == "PUT" {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
		checkErr(err)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Del("User-Agent")
		fmt.Println(req.Header.Get("User-Agent"))
		req.Header.Add("User-Agent", "Docker-Client/19.03.6")
		fmt.Println(req.Header.Get("User-Agent"))
		resp, err = client.Do(req)
		checkErr(err)

	} else {
		req, err := http.NewRequest(method, url, nil)
		checkErr(err)
		req.Header.Set("User-Agent", "Docker-Client/19.03.6 (linux)")
		resp, err = client.Do(req)
		checkErr(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	if resp.StatusCode > 250 { // 201 for container creation
		log.Error(string(body))
	}

	return resp
}

func other(w http.ResponseWriter, req *http.Request) {
	log.Info("New request received: " + req.Method + " " + req.URL.String())
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println("New request received: " + req.Method + " " + req.URL.String())
	fmt.Println(string(body))
	checkErr(err)
	serveReverseProxy(w, req, body)
}

func Start() *mux.Router {
	r := mux.NewRouter()

		r.HandleFunc("/containers/create", processCreateContainer)
		r.HandleFunc("/v{id}/containers/create", processCreateContainer)
		//r.HandleFunc("/v{id}/containers/{key}/wait", ProductHandler)
		r.PathPrefix("/").HandlerFunc(other)
		http.Handle("/", r)

		fmt.Println("Listening on port 8999")
		log.Fatal(http.ListenAndServe(":8999", nil))
		return r

	return r
}

// Serve a reverse proxy for a given url
func serveReverseProxy(res http.ResponseWriter, req *http.Request, body []byte) {
	// parse the url
	url, _ := url.Parse("http://localhost:2376")

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	//req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}
*/

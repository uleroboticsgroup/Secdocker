package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type plugin string

type config struct {
	Url         string `yaml:"url"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	AmountVulns int    `yaml:"amountvulns"`
}

type ImageStatus []struct {
	AnalysisStatus string `json:"analysis_status"`
}

type VulnAnalysis struct {
	Vulnerabilities []string `json:"vulnerabilities"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (p plugin) Process(image string) bool {
	log.Info("Anchore analysis initialized")

	f, err := os.Open("./plugins/anchore/config.yml")
	checkErr(err)
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	checkErr(err)

	//Generate auth header
	authData := "Basic " + base64.StdEncoding.EncodeToString([]byte(cfg.Username+":"+cfg.Password))

	// Add image
	url := cfg.Url + "/v1/images"
	var data = []byte(`{"tag":"` + image + `"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authData)
	_, err = client.Do(req)

	if err != nil && strings.Contains(err.Error(), "connection refused") {
		log.Error("Anchore is offline")
		return true
	} else {
		checkErr(err)
	}

	// Wait for validation of the image
	for {
		req, err := http.NewRequest("GET", cfg.Url+"/v1/images/by_id/"+image, nil)
		checkErr(err)
		req.Header.Add("Authorization", authData)
		resp, err := client.Do(req)
		checkErr(err)

		var imageStatus ImageStatus
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &imageStatus)
		checkErr(err)

		if imageStatus[0].AnalysisStatus == "analyzed" {
			break
		}

		time.Sleep(time.Second * 1)
	}

	// Get results
	req, err = http.NewRequest("GET", cfg.Url+"/v1/images/by_id/"+image+"+/vuln/all", nil)
	checkErr(err)
	req.Header.Add("Authorization", authData)
	resp, err := client.Do(req)
	checkErr(err)

	var imageVulns VulnAnalysis
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	err = json.Unmarshal(body, &imageVulns)

	log.Info("Vulnerabilities:")
	log.Info(imageVulns.Vulnerabilities)

	if len(imageVulns.Vulnerabilities) > cfg.AmountVulns {
		return false
	} else {
		return true
	}

}

var Plugin plugin

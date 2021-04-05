package main

import (
	"os"

	"github.com/uleroboticsgroup/Secdocker/tcpintercept"

	log "github.com/sirupsen/logrus"
)

func main() {
	file, err := os.OpenFile("secdocker.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	log.SetOutput(file)

	/// API SERVICE
	tcpintercept.ServeConnection()

}

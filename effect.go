package main

import (
	"log"
)

func main() {
	// Connect client
	lc, err := NewLightsClient("unix", "/tmp/monitor") //TODO
	if err != nil {
		log.Fatal(err)
	}
	defer lc.Close()
	// Start RPC serv
	// Start lighter
	result, err := lc.GetLightState("*")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result[0])
	log.Printf("%+v\n", result[0].HSBK)
}

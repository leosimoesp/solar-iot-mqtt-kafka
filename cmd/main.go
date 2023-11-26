package main

import (
	"github.com/lbsti/solar-iot-mqtt-kafka/core/usecase/simulator"
	"github.com/lbsti/solar-iot-mqtt-kafka/infra/envparser"
)

func main() {
	envParser := envparser.New()
	s := simulator.NewSimulator(envParser)
	s.Run()
}

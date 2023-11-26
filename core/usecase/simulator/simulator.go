package simulator

import (
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
)

type Simulator interface {
	Run()
}

type simulatorImpl struct {
	envParser dto.EnvParser
	cfg       *dto.Config
}

func NewSimulator(p dto.EnvParser) Simulator {
	return &simulatorImpl{envParser: p}
}

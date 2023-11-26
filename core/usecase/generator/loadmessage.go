package generator

import (
	"sync"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
)

func (sg solarGenerator) Load(chanMsgs chan<- []dto.Message, g *sync.WaitGroup, msgs []dto.Message) {
	numberOfSensors := sg.cfg.GetTotalOfSensors() * len(sg.cfg.GetPlants())
	for i := 0; i < len(msgs); i += numberOfSensors {
		start := i
		end := start + numberOfSensors
		if end >= len(msgs) {
			end = len(msgs)
		}
		chanMsgs <- msgs[start:end]
		g.Done()
	}
}

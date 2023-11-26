package config

import (
	"fmt"
	"log"
	"time"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
)

func Create(params dto.SolarSimulatorParams) dto.SolarSimulatorConfig {
	simulatorConfig := dto.NewSolarSimulatorConfig(params.TotalOfSensors, params.Plants)

	if params.TimeFailureStart != nil && params.TimeFailureEnd != nil {
		start, e := time.Parse("2006-01-02 15:04:05", *params.TimeFailureStart)
		handleErr("Failure when parse start date", e)
		end, e := time.Parse("2006-01-02 15:04:05", *params.TimeFailureEnd)
		handleErr("Failure when parse end date", e)

		if !start.IsZero() && !end.IsZero() {
			simulatorConfig.WithFailureAndRecover(start, end)
		}
	}
	simulatorConfig.WithSimulationIntervalInSecs(params.TimeSimulationIntervalSeconds)

	start, e := time.Parse("2006-01-02 15:04:05", params.TimeStartSimulation)
	handleErr("Failure when parse start simulation date", e)
	simulatorConfig.WithTimeStart(start)

	return simulatorConfig
}

func handleErr(msg string, e error) {
	if e != nil {
		log.Fatalf(fmt.Sprintf("%s %v", msg, e))
	}
}

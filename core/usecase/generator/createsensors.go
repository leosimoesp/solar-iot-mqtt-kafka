package generator

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/entity"
)

const sensorKeySize = 10

func (g solarGenerator) CreateSensors() ([]entity.SolarSensor, error) {
	plants := g.cfg.GetPlants()
	totalOfSensors := g.cfg.GetTotalOfSensors() * len(plants)
	createdSensors := make([]entity.SolarSensor, totalOfSensors)
	plantIndx := 0

	for i := 0; i < totalOfSensors; i++ {
		if plantIndx >= len(plants) {
			plantIndx = 0
		}
		createdSensor := createSensor(plants[plantIndx])

		if failureStart, failureEnd := g.cfg.GetTimeFailureAndRecover(); !failureStart.IsZero() &&
			!failureEnd.IsZero() && (i == 0 || i == 10) {
			createdSensor.ForceFailure()
		}

		createdSensors[i] = createdSensor
		plantIndx++
	}
	return createdSensors, nil
}

func createSensor(plantID int) entity.SolarSensor {
	builder := entity.NewSolarSensor()

	now := time.Now()
	timeFmt := now.Format("2006-01-02 15:04:05")

	return builder.
		WithPlantID(strconv.Itoa(plantID)).
		WithDateTime(timeFmt).
		WithSourceKey(generateSourceKey(sensorKeySize)).
		Build()
}

func generateSourceKey(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length+2)
	r.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

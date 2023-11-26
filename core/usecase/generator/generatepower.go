package generator

import (
	"math"
	"math/rand"
	"time"

	"github.com/lbsti/solar-iot-mqtt-kafka/adapter/mqtt"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/entity"
)

const (
	MinValuePower            = 0.0       //Watts
	MaxValuePower            = 10235.125 //Watts
	TotalOfSolarGenInSeconds = 44100
)

func (sg solarGenerator) GeneratePower(sensors []entity.SolarSensor) ([]dto.Message, error) {
	start := sg.cfg.GetTimeStart()
	end := start.Add(time.Second * TotalOfSolarGenInSeconds)
	diffInSecs := int(end.Sub(start).Seconds())
	intervalInSecs := sg.cfg.GetSimulationIntervalInSecs()
	timeSlots := diffInSecs / intervalInSecs
	totalOfMessages := len(sensors) * timeSlots
	messages := make([]dto.Message, 0, totalOfMessages)
	indxSensors := 0
	now := start
	_, failureEnd := sg.cfg.GetTimeFailureAndRecover()

	for i := 0; i < totalOfMessages; i++ {
		solarSensor := sensors[indxSensors]
		solarSensor.RunFailureOrRecover(now, failureEnd)
		acPower, dcPower := generateSolarPower(now)

		if solarSensor.IsFailure() {
			solarSensor.WithPower(0.0, 0.0)
		} else {
			solarSensor.WithPower(acPower, dcPower)
		}
		solarSensor.WithDateTime(now)
		messages = append(messages, mqtt.Create(solarSensor))
		if indxSensors == len(sensors)-1 {
			now = now.Add(time.Second * time.Duration(intervalInSecs))
		}
		indxSensors++
		if indxSensors >= len(sensors) {
			indxSensors = 0
		}
	}
	return messages, nil
}

func generateSolarPower(currentTime time.Time) (float64, float64) {

	factor := 0.0

	switch currentTime.Hour() {
	case 0, 1, 2, 3, 4, 5, 19, 20, 21, 22, 23, 24:
		factor = 0.0
	case 6, 7, 17, 18:
		factor = 0.1
	case 8, 9:
		factor = 0.4
	case 10, 11, 15, 16:
		factor = 0.9
	case 12, 13, 14:
		factor = 1.4
	}
	dcPower := generateRandFloat(MinValuePower, MaxValuePower) * factor
	return math.Floor(dcPower*100) / 100, calculateACPower(dcPower)
}

func generateRandFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func calculateACPower(dcPower float64) float64 {
	acPower := dcPower * 0.09
	return math.Floor(acPower*100) / 100
}

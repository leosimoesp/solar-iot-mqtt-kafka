package generator

import (
	"testing"
	"time"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/stretchr/testify/assert"
)

func Test_solarGenerator_CreateSensors(t *testing.T) {
	t.Run("Should create sensors successfully when there isn't a failure", createSolarSensorsWithoutFail)
	t.Run("Should create sensors successfully when there is a failure", createSolarSensorsWithFail)
}

func createSolarSensorsWithoutFail(t *testing.T) {
	totalOfSensors := 10
	plants := []int{1, 2}

	cfg := dto.NewSolarSimulatorConfig(totalOfSensors, plants)

	g := NewSolarGenerator(cfg)
	sensors, e := g.CreateSensors()
	assert.Nil(t, e)
	assert.NotEmpty(t, sensors)
	assert.Equal(t, 20, len(sensors))
	assert.Equal(t, false, sensors[0].IsFailure())
}

func createSolarSensorsWithFail(t *testing.T) {
	totalOfSensors := 10
	plants := []int{1, 2}
	cfg := dto.NewSolarSimulatorConfig(totalOfSensors, plants)
	start := time.Now().Add(time.Minute * 2)
	end := start.Add(time.Minute * 25)
	cfg.WithFailureAndRecover(start, end)

	g := NewSolarGenerator(cfg)
	sensors, e := g.CreateSensors()
	assert.Nil(t, e)
	assert.NotEmpty(t, sensors)
	assert.Equal(t, 20, len(sensors))
	assert.Equal(t, true, sensors[0].IsFailure())
	assert.Equal(t, false, sensors[1].IsFailure())
	assert.Equal(t, true, sensors[10].IsFailure())
}

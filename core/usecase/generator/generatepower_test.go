package generator

import (
	"testing"
	"time"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/stretchr/testify/assert"
)

func Test_solarGenerator_GeneratePower(t *testing.T) {
	t.Run("Should simulate sensor solar power generation without failure", simulationWithoutFailure)
	t.Run("Should simulate sensor solar power generation with failure", simulationWithFailure)
	t.Run("Should simulate sensor solar power recover after time failure expires", simulationFailureRecover)
}

func simulationWithoutFailure(t *testing.T) {
	totalOfSensors := 10
	plants := []int{1, 2}
	cfg := dto.NewSolarSimulatorConfig(totalOfSensors, plants)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.UTC)

	cfg.WithTimeStart(start)
	cfg.WithSimulationIntervalInSecs(900)

	s := NewSolarGenerator(cfg)
	sensors, e := s.CreateSensors()
	assert.Nil(t, e)
	assert.NotEmpty(t, sensors)
	assert.Equal(t, 20, len(sensors))
	assert.Equal(t, false, sensors[0].IsFailure())

	messages, e := s.GeneratePower(sensors)
	assert.Nil(t, e)
	assert.NotEmpty(t, messages)
	assert.Equal(t, 980, len(messages))
	assert.NotNil(t, messages[0])
}

func simulationWithFailure(t *testing.T) {
	totalOfSensors := 10
	plants := []int{1, 2}
	cfg := dto.NewSolarSimulatorConfig(totalOfSensors, plants)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.UTC)
	failureStart, failureEnd := start.Add(time.Second*5), start.Add(time.Hour*2)

	cfg.WithTimeStart(start)
	cfg.WithSimulationIntervalInSecs(900)
	cfg.WithFailureAndRecover(failureStart, failureEnd)

	s := NewSolarGenerator(cfg)
	sensors, e := s.CreateSensors()
	assert.Nil(t, e)
	assert.NotEmpty(t, sensors)
	assert.Equal(t, 20, len(sensors))
	assert.Equal(t, true, sensors[0].IsFailure())

	messages, e := s.GeneratePower(sensors)
	assert.Nil(t, e)
	assert.NotEmpty(t, messages)
	assert.Equal(t, 980, len(messages))
	assert.NotNil(t, messages[0])
}

func simulationFailureRecover(t *testing.T) {
	totalOfSensors := 10
	plants := []int{1, 2}
	cfg := dto.NewSolarSimulatorConfig(totalOfSensors, plants)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.UTC)
	failureStart, failureEnd := start.Add(time.Second*150), start.Add(time.Hour*2)

	cfg.WithTimeStart(start)
	cfg.WithSimulationIntervalInSecs(900)
	cfg.WithFailureAndRecover(failureStart, failureEnd)

	s := NewSolarGenerator(cfg)
	sensors, e := s.CreateSensors()
	assert.Nil(t, e)
	assert.NotEmpty(t, sensors)
	assert.Equal(t, 20, len(sensors))
	assert.Equal(t, true, sensors[0].IsFailure())

	if start.Before(failureStart) {
		assert.Equal(t, true, sensors[0].IsFailure())
	} else {
		assert.Equal(t, false, sensors[0].IsFailure())
	}

	messages, e := s.GeneratePower(sensors)
	assert.Nil(t, e)
	assert.NotEmpty(t, messages)
	assert.Equal(t, 980, len(messages))
	assert.NotNil(t, messages[0])

}

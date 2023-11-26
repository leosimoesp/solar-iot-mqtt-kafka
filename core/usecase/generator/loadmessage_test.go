package generator

import (
	"sync"
	"testing"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/stretchr/testify/assert"
)

func Test_solarGenerator_Load(t *testing.T) {
	totalOfSensors := 2
	plants := []int{1, 2}
	cfg := dto.NewSolarSimulatorConfig(totalOfSensors, plants)

	msgs := []dto.Message{
		{
			Payload: `{"id":"1"}`,
		},
		{
			Payload: `{"id":"2"}`,
		},
		{
			Payload: `{"id":"3"}`,
		},
		{
			Payload: `{"id":"4"}`,
		},
	}

	gen := NewSolarGenerator(cfg)

	chanMsgs := make(chan []dto.Message)
	var wg sync.WaitGroup

	wg.Add(2)
	go gen.Load(chanMsgs, &wg, msgs)

	go func(d *sync.WaitGroup) {
		msgs := <-chanMsgs
		assert.Equal(t, 4, len(msgs))
		d.Done()
	}(&wg)

	go func(d *sync.WaitGroup) {
		d.Wait()
		defer close(chanMsgs)
	}(&wg)
}

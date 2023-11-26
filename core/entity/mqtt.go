package entity

import (
	"context"
	"sync"

	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
)

type MqttProducer interface {
	Produce(ctx context.Context, g *sync.WaitGroup, sensorsData <-chan []dto.Message)
}

type MqttMessageLoader interface {
	Load(chanMsgs chan<- []dto.Message, g *sync.WaitGroup, msgs []dto.Message)
}

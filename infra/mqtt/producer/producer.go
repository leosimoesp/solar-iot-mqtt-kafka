package producer

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/eclipse/paho.golang/paho"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/entity"
)

type mqttProducer struct {
	mqttClient *paho.Client
	cfg        *dto.MqttServerConfig
}

func NewMqttProducer(cfg dto.MqttServerConfig, client *paho.Client) entity.MqttProducer {
	p := &mqttProducer{
		mqttClient: client,
		cfg:        &cfg,
	}
	return p
}

// Produce implements mqtt.MqttProducer.
func (p *mqttProducer) Produce(ctx context.Context, g *sync.WaitGroup, sensorsData <-chan []dto.Message) {
	for {
		select {
		case <-ctx.Done():
			g.Done()
			return
		case data := <-sensorsData:
			logMessage(data)
			log.Default().Println("Simulating wait time to next sensors data collect event.......")
			time.Sleep(10 * time.Second)
			p.send(ctx, data)
			g.Done()
		}
	}
}

// send uses the mqtt client connected and send the messages to mosquitto
func (p *mqttProducer) send(ctx context.Context, messages []dto.Message) {
	var wg sync.WaitGroup
	for _, message := range messages {
		wg.Add(1)
		go func(g *sync.WaitGroup, msg dto.Message) {
			defer g.Done()
			props := &paho.PublishProperties{}
			props.User.Add("clientId", p.mqttClient.ClientID)

			pb := &paho.Publish{
				Topic:      p.cfg.MqttTopic,
				QoS:        byte(0),
				Payload:    []byte(msg.Payload.(string)),
				Properties: props,
			}

			if _, err := p.mqttClient.Publish(context.Background(), pb); err != nil {
				log.Fatalf("Error when publish%v\n", err)
			}
		}(&wg, message)
	}
	wg.Wait()
}

// logMessage show at stdio the messages with line separator
func logMessage(messages []dto.Message) {
	repeated := strings.Repeat(" ", 110)
	lineSeparator := strings.Repeat("-", utf8.RuneCountInString(repeated))
	for _, sensor := range messages {
		log.Default().Println(sensor)
	}
	log.Default().Println(lineSeparator)
}

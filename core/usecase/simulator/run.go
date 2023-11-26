package simulator

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/eclipse/paho.golang/paho"
	"github.com/lbsti/solar-iot-mqtt-kafka/adapter/config"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/entity"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/usecase/generator"
	"github.com/lbsti/solar-iot-mqtt-kafka/infra/mqtt/client/pahomqtt"
	"github.com/lbsti/solar-iot-mqtt-kafka/infra/mqtt/producer"
)

// Run implements Simulator.
func (s *simulatorImpl) Run() {
	s.cfg = s.loadEnvs()
	powerGenerator := generator.NewSolarGenerator(s.cfg.SimulatorCfg)
	sensors, e := powerGenerator.CreateSensors()
	if e != nil {
		log.Fatalf("Failed to create sensors: %v", e)
	}
	ctx := context.Background()
	mqttClient := pahomqtt.NewClient(s.cfg.MqttServerCfg.BrokerServerUrl, paho.ClientConfig{})
	mqttProducer := s.createMqttProducer(ctx, mqttClient)

	mqttMsgs, e := powerGenerator.GeneratePower(sensors)
	if e != nil {
		log.Fatalf("Failed to generate random sensors: %v", e)
	}

	var wg sync.WaitGroup
	chanSensors := make(chan []dto.Message)
	delta := 2 * (len(mqttMsgs) / len(sensors))
	wg.Add(delta)

	go powerGenerator.Load(chanSensors, &wg, mqttMsgs)
	go mqttProducer.Produce(ctx, &wg, chanSensors)
	wg.Wait()

	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ic
		d := &paho.Disconnect{ReasonCode: 0}
		err := mqttClient.Disconnect(d)
		if err != nil {
			log.Fatalf("Failed to send Disconnect: %v", err)
		}
	}()
}

func (s simulatorImpl) loadEnvs() *dto.Config {
	mqttServerCfg := dto.MqttServerConfig{}
	if e := s.envParser.Parse(&mqttServerCfg); e != nil {
		log.Fatalf("unable to parse .env file: %e", e)
	}
	simulatorParams := dto.SolarSimulatorParams{}
	if e := s.envParser.Parse(&simulatorParams); e != nil {
		log.Fatalf("unable to parse .env file: %e", e)
	}
	simulatorCfg := config.Create(simulatorParams)
	cfg := dto.Config{}
	cfg.MqttServerCfg = mqttServerCfg
	cfg.SimulatorCfg = simulatorCfg
	return &cfg
}

func (s simulatorImpl) createMqttProducer(ctx context.Context, mqttClient *paho.Client) entity.MqttProducer {
	mqttConnectCfg := pahomqtt.Connect(dto.MqttClientConfig{
		ClientID:     s.cfg.MqttServerCfg.MqttClientID,
		Username:     s.cfg.MqttServerCfg.MqttUsername,
		Password:     s.cfg.MqttServerCfg.MqttPassword,
		KeepAlive:    30,
		CleanStart:   true,
		UsernameFlag: true,
		PasswordFlag: true,
	})

	if _, e := mqttClient.Connect(ctx, mqttConnectCfg); e != nil {
		log.Fatalf("Failed to connect to server: %v", e)
	}

	return producer.NewMqttProducer(s.cfg.MqttServerCfg, mqttClient)
}

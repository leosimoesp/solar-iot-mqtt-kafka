package pahomqtt

import (
	"log"
	"net"

	mqtt "github.com/eclipse/paho.golang/paho"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
)

func NewClient(server string, cfg mqtt.ClientConfig) *mqtt.Client {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", server, err)
	}

	cfg.Router = mqtt.NewSingleHandlerRouter(func(m *mqtt.Publish) {
		log.Printf("%v : %v %v", m.Topic, m.Properties.User.Get("clientId"), string(m.Payload))
	})
	cfg.Conn = conn
	return mqtt.NewClient(cfg)
}

func Connect(cfg dto.MqttClientConfig) *mqtt.Connect {
	cp := &mqtt.Connect{
		KeepAlive:    uint16(cfg.KeepAlive),
		ClientID:     cfg.ClientID,
		CleanStart:   cfg.CleanStart,
		Username:     cfg.Username,
		Password:     []byte(cfg.Password),
		UsernameFlag: cfg.UsernameFlag,
		PasswordFlag: cfg.PasswordFlag,
	}
	return cp
}

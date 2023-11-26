package mqtt

import (
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/entity"
)

func Create(sensor entity.SolarSensor) dto.Message {
	return dto.Message{
		Payload: sensor.String(),
	}
}

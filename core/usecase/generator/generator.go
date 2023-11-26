package generator

import (
	"github.com/lbsti/solar-iot-mqtt-kafka/core/dto"
	"github.com/lbsti/solar-iot-mqtt-kafka/core/entity"
)

type SolarPowerGenerator interface {
	CreateSensors() ([]entity.SolarSensor, error)
	GeneratePower(sensors []entity.SolarSensor) ([]dto.Message, error)
	entity.MqttMessageLoader
}

type solarGenerator struct {
	cfg dto.SolarSimulatorConfig
}

func NewSolarGenerator(c dto.SolarSimulatorConfig) SolarPowerGenerator {
	return &solarGenerator{
		cfg: c,
	}
}

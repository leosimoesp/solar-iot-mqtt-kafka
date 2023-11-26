package dto

import "time"

type SolarSimulatorConfig interface {
	GetTotalOfSensors() int
	GetTimeFailureAndRecover() (start, end time.Time)
	WithFailureAndRecover(start, end time.Time)
	GetPlants() []int
	GetTimeStart() (start time.Time)
	GetSimulationIntervalInSecs() int
	WithSimulationIntervalInSecs(interval int)
	WithTimeStart(t time.Time)
}

type solarSimulatorConfig struct {
	totalOfSensors            int
	plants                    []int
	timeStartSimulation       time.Time
	simulationIntervalSeconds int
	timeFailureStart          *time.Time
	timeFailureEnd            *time.Time
}

// WithSimulationIntervalInSecs implements SolarSimulatorConfig.
func (cfg *solarSimulatorConfig) WithSimulationIntervalInSecs(interval int) {
	cfg.simulationIntervalSeconds = interval
}

// WithTimeStart implements SolarSimulatorConfig.
func (cfg *solarSimulatorConfig) WithTimeStart(t time.Time) {
	cfg.timeStartSimulation = t
}

// GetSimulationIntervalInSecs implements SolarSimulatorConfig.
func (cfg solarSimulatorConfig) GetSimulationIntervalInSecs() int {
	return cfg.simulationIntervalSeconds
}

// GetTimeStart implements SolarSimulatorConfig.
func (cfg solarSimulatorConfig) GetTimeStart() (start time.Time) {
	return cfg.timeStartSimulation
}

// GetPlants implements SolarSimulatorConfig.
func (cfg solarSimulatorConfig) GetPlants() []int {
	return cfg.plants
}

// GetTimeFailureAndRecover implements SolarSimulatorConfig.
func (cfg *solarSimulatorConfig) GetTimeFailureAndRecover() (start time.Time, end time.Time) {

	if cfg.timeFailureStart != nil {
		start = *cfg.timeFailureStart
	}
	if cfg.timeFailureEnd != nil {
		end = *cfg.timeFailureEnd
	}
	return start, end
}

// GetTotalOfSensors implements SolarSimulatorConfig.
func (cfg solarSimulatorConfig) GetTotalOfSensors() int {
	return cfg.totalOfSensors
}

// WithFailureAndRecover implements SolarSimulatorConfig.
func (cfg *solarSimulatorConfig) WithFailureAndRecover(start time.Time, end time.Time) {
	cfg.timeFailureStart = &start
	cfg.timeFailureEnd = &end
}

func NewSolarSimulatorConfig(totalOfSensors int, plants []int) SolarSimulatorConfig {
	return &solarSimulatorConfig{
		totalOfSensors:   totalOfSensors,
		plants:           plants,
		timeFailureStart: nil,
		timeFailureEnd:   nil,
	}
}

type SolarSimulatorParams struct {
	TotalOfSensors                int     `env:"TOTAL_OF_SENSORS,required"`
	Plants                        []int   `env:"PLANTS,required"`
	TimeStartSimulation           string  `env:"TIME_START_SIMULATION,required"`
	TimeSimulationIntervalSeconds int     `env:"TIME_SIMULATION_INTERVAL_SECONDS,required"`
	TimeFailureStart              *string `env:"TIME_FAILURE_START"`
	TimeFailureEnd                *string `env:"TIME_FAILURE_END"`
}

type MqttServerConfig struct {
	BrokerServerUrl string `env:"BROKER_SERVER_URL,required"`
	MqttTopic       string `env:"MQTT_TOPIC,required"`
	MqttClientID    string `env:"MQTT_CLIENTID,required"`
	MqttPassword    string `env:"MQTT_PASSWORD,required"`
	MqttUsername    string `env:"MQTT_USERNAME,required"`
}

type MqttClientConfig struct {
	ClientID     string
	Username     string
	Password     string
	KeepAlive    int
	CleanStart   bool
	UsernameFlag bool
	PasswordFlag bool
}

type Config struct {
	SimulatorCfg  SolarSimulatorConfig
	MqttServerCfg MqttServerConfig
}

type EnvParser interface {
	Parse(v interface{}) error
}

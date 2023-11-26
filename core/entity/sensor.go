package entity

import (
	"encoding/json"
	"log"
	"time"
)

type SolarSensor interface {
	WithPower(acPower float64, dcPower float64)
	ForceFailure()
	IsFailure() bool
	String() string
	WithDateTime(t time.Time)
	RunFailureOrRecover(current, end time.Time)
}

type sensor struct {
	PlantID   string  `json:"plantId"`
	SourceKey string  `json:"sourceKey"`
	DateTime  string  `json:"datetime"`
	DCPower   float64 `json:"dcPower"`
	ACPower   float64 `json:"acPower"`
	failure   bool
}

// RunFailureOrRecover implements SolarSensor.
func (s *sensor) RunFailureOrRecover(current time.Time, end time.Time) {
	if s.failure && !end.IsZero() && current.After(end) {
		s.failure = false
	}
}

// WithDateTime implements SolarSensor.
func (s *sensor) WithDateTime(dateTime time.Time) {
	timeFmt := dateTime.Format("2006-01-02 15:04:05")
	s.DateTime = timeFmt
}

// String implements SolarSensor.
func (s *sensor) String() string {
	resultAsStr, e := json.Marshal(s)

	if e != nil {
		log.Fatalf("Error %v", e)
		return e.Error()
	}
	return string(resultAsStr)
}

// IsFailure implements SolarSensor.
func (s sensor) IsFailure() bool {
	return s.failure
}

// ForceFailure implements SolarSensor.
func (s *sensor) ForceFailure() {
	s.failure = true
}

// WithPower implements SolarSensor.
func (s *sensor) WithPower(acPower float64, dcPower float64) {
	s.ACPower = acPower
	s.DCPower = dcPower
}

type SensorBuilder interface {
	WithPlantID(string) SensorBuilder
	WithSourceKey(string) SensorBuilder
	WithDateTime(string) SensorBuilder
	WithDCPower(float64) SensorBuilder
	WithACPower(float64) SensorBuilder
	WithFailure() SensorBuilder
	Build() SolarSensor
}

type sensorBuilder struct {
	sensor sensor
}

// WithFailure implements SensorBuilder.
func (sb *sensorBuilder) WithFailure() SensorBuilder {
	sb.sensor.failure = true
	return sb
}

// Build implements SensorBuilder.
func (sb *sensorBuilder) Build() SolarSensor {
	return &sb.sensor
}

// WithACPower implements SensorBuilder.
func (sb *sensorBuilder) WithACPower(acPower float64) SensorBuilder {
	sb.sensor.ACPower = acPower
	return sb
}

// WithDCPower implements SensorBuilder.
func (sb *sensorBuilder) WithDCPower(dCPower float64) SensorBuilder {
	sb.sensor.DCPower = dCPower
	return sb
}

// WithDateTime implements SensorBuilder.
func (sb *sensorBuilder) WithDateTime(dateTime string) SensorBuilder {
	sb.sensor.DateTime = dateTime
	return sb
}

// WithPlantID implements SensorBuilder.
func (sb *sensorBuilder) WithPlantID(plantID string) SensorBuilder {
	sb.sensor.PlantID = plantID
	return sb
}

// WithSourceKey implements SensorBuilder.
func (sb *sensorBuilder) WithSourceKey(sourceKey string) SensorBuilder {
	sb.sensor.SourceKey = sourceKey
	return sb
}

func NewSolarSensor() SensorBuilder {
	return &sensorBuilder{
		sensor: sensor{},
	}
}

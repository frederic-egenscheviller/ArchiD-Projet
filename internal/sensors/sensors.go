package sensors

import (
	"ArchiD-Projet/airportsensors/meteofranceAPI"
	"ArchiD-Projet/internal/mqttconnect"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type SensorData struct {
	SensorID         int
	AirportID        string
	Measurement      string
	MeasurementValue float64
	MeasurementTime  time.Time
}

type SensorConfig struct {
	BrokerAddress         string        `yaml:"brokerAddress"`
	Port                  int           `yaml:"port"`
	QoS                   byte          `yaml:"qos"`
	ClientID              string        `yaml:"clientID"`
	TransmissionFrequency time.Duration `yaml:"transmissionFrequency"`
}

type Sensor struct {
	client    *mqttconnect.Client
	qos       byte
	topic     string
	retained  bool
	lastValue string
	Config    SensorConfig
}

func LoadSensorConfig(filename string) (SensorConfig, error) {
	var config SensorConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func NewSensor(client *mqttconnect.Client, qos byte, retained bool, config SensorConfig) *Sensor {
	return &Sensor{
		client:   client,
		qos:      qos,
		retained: retained,
		Config:   config,
	}
}

func (sensor *Sensor) PublishSensorData(data SensorData) {
	payload := fmt.Sprintf(`"%s %s %f"`,
		data.MeasurementTime.Format("2006-01-02 15:04:05"), data.Measurement, data.MeasurementValue)

	sensor.topic = "airport/" + data.AirportID

	err := sensor.client.Publish(sensor.topic, sensor.qos, sensor.retained, payload)
	if err != nil {
		return
	}
}

func (sensor *Sensor) StartMonitoring() {
	ticker := time.NewTicker(sensor.Config.TransmissionFrequency)
	go func() {
		for {
			select {
			case <-ticker.C:
				sensorData, err := meteofranceAPI.FetchSensorDataFromAPI(sensor.Config.ClientID)
				if err != nil {
					fmt.Println("Error fetching sensor data from API:", err)
					continue
				}
				sensor.PublishSensorData(SensorData(sensorData))
			}
		}
	}()

	mqttconnect.WaitForSignal()
}

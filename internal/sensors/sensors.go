package sensors

import (
	"ArchiD-Projet/internal/meteofranceAPI"
	"ArchiD-Projet/internal/mqttconnect"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

type Sensor struct {
	client    *mqttconnect.Client
	qos       byte
	topic     string
	retained  bool
	config    SensorConfig
	info      SensorInfo
	waitGroup *sync.WaitGroup
}

type SensorConfig struct {
	BrokerAddress         string        `yaml:"brokerAddress"`
	Port                  int           `yaml:"port"`
	QoS                   byte          `yaml:"qos"`
	ClientID              string        `yaml:"clientID"`
	TransmissionFrequency time.Duration `yaml:"transmissionFrequency"`
}

type SensorData struct {
	SensorID         int
	AirportID        string
	Measurement      string
	MeasurementValue float64
	MeasurementTime  time.Time
}

type RetrievedSensorsConfig struct {
	BrokerAddress string       `yaml:"brokerAddress"`
	Port          int          `yaml:"port"`
	Sensors       []SensorInfo `yaml:"sensors"`
}

type SensorInfo struct {
	TransmissionFrequency time.Duration `yaml:"transmissionFrequency"`
	ClientID              string        `yaml:"clientID"`
	QoS                   byte          `yaml:"qos"`
	AirportIATA           string        `yaml:"airportIATA"`
	GeoIDInsee            string        `yaml:"geoIDInsee"`
}

func LoadSensorConfigs(filename string) (RetrievedSensorsConfig, error) {
	var configs RetrievedSensorsConfig

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(configs)
		return configs, err
	}

	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return configs, err
	}

	return configs, nil
}

func NewSensor(client *mqttconnect.Client, qos byte, retained bool, config SensorConfig, info SensorInfo) *Sensor {
	return &Sensor{
		client:   client,
		qos:      qos,
		retained: retained,
		config:   config,
		info:     info,
	}
}

func (sensor *Sensor) PublishSensorData(data SensorData) {
	payload := fmt.Sprintf("%s %s %f",
		data.MeasurementTime.Format("2006-01-02 15:04:05"), data.Measurement, data.MeasurementValue)

	sensor.topic = "airports/" + data.AirportID

	err := sensor.client.Publish(sensor.topic, sensor.qos, sensor.retained, payload)
	if err != nil {
		return
	}
}

func (sensor *Sensor) StartMonitoring() {
	ticker := time.NewTicker(sensor.config.TransmissionFrequency)

	sensor.waitGroup.Add(1)

	go func() {
		defer sensor.waitGroup.Done()
		defer ticker.Stop()
		for range ticker.C {
			sensorData, err := meteofranceAPI.FetchSensorDataFromAPI(meteofranceAPI.SensorInfo(sensor.info))
			if err != nil {
				fmt.Println("Error fetching sensor data from API:", err)
				continue
			}
			sensor.PublishSensorData(SensorData(sensorData))
		}
	}()

	mqttconnect.WaitForSignal()
}

func LoadSensors(retrievedSensorsConfig RetrievedSensorsConfig) {
	var sensorsList []*Sensor

	for _, sensorInfo := range retrievedSensorsConfig.Sensors {
		client, err := mqttconnect.NewClient(retrievedSensorsConfig.BrokerAddress, sensorInfo.ClientID, nil)
		if err != nil {
			fmt.Printf("Error creating MQTT client for %s: %v\n", sensorInfo.ClientID, err)
			continue
		}

		config := SensorConfig{
			BrokerAddress:         retrievedSensorsConfig.BrokerAddress,
			Port:                  retrievedSensorsConfig.Port,
			QoS:                   sensorInfo.QoS,
			ClientID:              sensorInfo.ClientID,
			TransmissionFrequency: sensorInfo.TransmissionFrequency,
		}

		sensor := NewSensor(client, sensorInfo.QoS, true, config, sensorInfo)
		sensorsList = append(sensorsList, sensor)
	}

	var waitGroup sync.WaitGroup

	for _, sensor := range sensorsList {
		sensor.waitGroup = &waitGroup
		go sensor.StartMonitoring()
	}

	waitGroup.Wait()
	mqttconnect.WaitForSignal()
}

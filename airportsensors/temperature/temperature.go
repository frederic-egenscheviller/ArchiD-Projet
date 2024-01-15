package main

import (
	"ArchiD-Projet/internal/mqttconnect"
	"ArchiD-Projet/internal/sensors"
	"fmt"
)

func main() {
	sensorConfig, _ := sensors.LoadSensorConfig("airportsensors/temperature/temperature_sensor_config.yml")

	client, err := mqttconnect.NewClient(sensorConfig.BrokerAddress, sensorConfig.ClientID)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	defer client.Disconnect()

	sensor := sensors.NewSensor(client, sensorConfig.QoS, true, sensorConfig)

	sensor.StartMonitoring()
}

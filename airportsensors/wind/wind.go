package main

import (
	"ArchiD-Projet/airportsensors/meteofranceAPI"
	"ArchiD-Projet/internal/mqttconnect"
	"ArchiD-Projet/internal/sensors"
	"fmt"
	"time"
)

func main() {
	sensorConfig, _ := sensors.LoadSensorConfig("airportsensors/wind/wind_sensor_config.yml")

	client, err := mqttconnect.NewClient(sensorConfig.BrokerAddress, sensorConfig.ClientID)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	defer client.Disconnect()

	sensor := sensors.NewSensor(client, sensorConfig.QoS, true, sensorConfig)

	ticker := time.NewTicker(sensor.Config.TransmissionFrequency)
	go func() {
		for {
			select {
			case <-ticker.C:
				sensorData, err := meteofranceAPI.FetchSensorDataFromAPI(sensorConfig.ClientID)
				if err != nil {
					fmt.Println("Error fetching sensor data from API:", err)
					continue
				}
				sensor.PublishSensorData(sensorData)
			}
		}
	}()

	mqttconnect.WaitForSignal()
}

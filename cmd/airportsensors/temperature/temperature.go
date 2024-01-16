package main

import (
	"ArchiD-Projet/internal/sensors"
	"fmt"
)

func main() {
	retrievedSensorsConfig, err := sensors.LoadSensorConfigs("config/temperature_sensor_config.yml")
	if err != nil {
		fmt.Println("Error loading sensor configurations:", err)
		return
	}

	sensors.LoadSensors(retrievedSensorsConfig)
}

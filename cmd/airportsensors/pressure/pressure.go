package main

import (
	"ArchiD-Projet/internal/sensors"
	"log"
)

func main() {
	retrievedSensorsConfig, err := sensors.LoadSensorConfigs("config/pressure_sensor_config.yml")
	if err != nil {
		log.Fatal("Error loading sensor configurations:", err)
		return
	}

	sensors.LoadSensors(retrievedSensorsConfig)
}

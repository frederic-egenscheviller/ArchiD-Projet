package main

import (
	brokerConfiguration "ArchiD-Projet/internal/brokerConfiguration"
	brokerUtils "ArchiD-Projet/internal/brokerUtils"
	"ArchiD-Projet/internal/mqttconnect"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strings"
	"time"
)

var (
	config = brokerConfiguration.GetFileRecorderSettings()
	BROKER = brokerConfiguration.GetBrokerAddress()
	TOPIC  = config[0]
)

func onMessageReceived(_ mqtt.Client, message mqtt.Message) {
	payload := string(message.Payload())
	data := strings.Split(payload, " ")

	submittedTimestamp := data[0] + "T" + data[1] + "Z"
	sensor := data[2]
	value := data[3]

	timestamp, err := time.Parse(time.RFC3339, submittedTimestamp)
	if err != nil {
		log.Println("Failed to parse timestamp:", err)
		return
	}

	fileName := fmt.Sprintf("%s_%s.csv", brokerUtils.GetAirportCodeFromTopic(message.Topic()), timestamp.Format("2006-01-02"))

	if _, err := os.Stat(config[1]); os.IsNotExist(err) {
		err := os.Mkdir(config[1], 0755)
		if err != nil {
			log.Fatal("Failed to create folder:", err)
			return
		}
	}

	file, err := os.OpenFile(config[1]+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open or create file:", err)
		return
	}
	defer file.Close()

	line := fmt.Sprintf("%s %s %s\n", timestamp.Format("2006-01-02 15:04:05"), sensor, value)
	_, err = file.WriteString(line)
	if err != nil {
		log.Println("Failed to write to file:", err)
		return
	}
}

func main() {
	client, err := mqttconnect.NewClient(BROKER, "file_recorder", onMessageReceived)
	if err != nil {
		log.Fatal("Error creating MQTT client:", err)
		return
	}

	err = client.Subscribe(TOPIC, 1, nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic:", err)
		return
	}

	mqttconnect.WaitForSignal()
}

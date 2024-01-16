package main

import (
	brokerconfiguration "ArchiD-Projet/internal/brokerConfiguration"
	brokerutils "ArchiD-Projet/internal/brokerUtils"
	"ArchiD-Projet/internal/mqttconnect"
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

var config = brokerconfiguration.GetDatabaseRecorderSettings()

var (
	BROKER = brokerconfiguration.GetBrokerAddress()
	BUCKET = config[0]
	ORG    = config[1]
	URL    = config[2]
	TOPIC  = config[3]
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Println(client.IsConnected())
	payload := string(message.Payload())
	data := strings.Split(payload, " ")

	apiKey := os.Getenv("INFLUX_DB_API_KEY")

	influxClient := influxdb2.NewClient(URL, apiKey)

	submittedTimestamp := data[0] + " " + data[1]
	sensor := data[2]
	value, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		fmt.Println("Failed to convert value to integer")
		return
	}

	parsedTime, err := time.Parse(time.DateTime, submittedTimestamp)
	if err != nil {
		fmt.Println("Erreur de conversion de la chaîne en objet time.Time :", err)
		return
	}

	writeAPI := influxClient.WriteAPIBlocking(ORG, BUCKET)

	fmt.Println(parsedTime)
	p := influxdb2.NewPoint(brokerutils.GetAirportCodeFromTopic(message.Topic()),
		map[string]string{"unit": sensor},
		map[string]interface{}{"value": value},
		time.Now())

	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return
	}

	influxClient.Close()
}

func main() {
	client, err := mqttconnect.NewClient(BROKER, "database_recorder", onMessageReceived)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	err = client.Subscribe(TOPIC, 1, nil)
	if err != nil {
		return
	}
	mqttconnect.WaitForSignal()
}

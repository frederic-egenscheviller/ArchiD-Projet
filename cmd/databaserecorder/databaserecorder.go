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

var (
	config       = brokerconfiguration.GetDatabaseRecorderSettings()
	BROKER       = brokerconfiguration.GetBrokerAddress()
	BUCKET       = config[0]
	ORG          = config[1]
	URL          = config[2]
	TOPIC        = config[3]
	influxClient influxdb2.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func onMessageReceived(_ mqtt.Client, message mqtt.Message) {

	payload := string(message.Payload())
	data := strings.Split(payload, " ")

	submittedTimestamp := data[0] + " " + data[1]
	sensor := data[2]
	value, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		fmt.Println("Failed to convert value to integer")
		return
	}

	loc, _ := time.LoadLocation("Europe/Paris")
	th, _ := time.ParseInLocation("02/01/2006 15:04:05", submittedTimestamp, loc)

	writeAPI := influxClient.WriteAPIBlocking(ORG, BUCKET)

	p := influxdb2.NewPointWithMeasurement(sensor).
		AddTag("airport", brokerutils.GetAirportCodeFromTopic(message.Topic())).
		AddField("value", value).
		SetTime(th)

	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return
	}
}

func main() {
	apiKey := os.Getenv("INFLUX_DB_API_KEY")
	if apiKey == "" {
		fmt.Println("INFLUX_DB_API_KEY environment variable not set")
		return
	}
	influxClient = influxdb2.NewClient(URL, apiKey)
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

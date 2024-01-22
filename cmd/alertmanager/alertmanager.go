package main

import (
	"ArchiD-Projet/internal/brokerconfiguration"
	"ArchiD-Projet/internal/brokerutils"
	"ArchiD-Projet/internal/mqttconnect"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Thresholds struct {
	Temp struct {
		Min float64 `yaml:"min"`
		Max float64 `yaml:"max"`
	} `yaml:"temp"`
	Wind struct {
		Speed float64 `yaml:"speed"`
	} `yaml:"wind"`
	Pressure struct {
		Summer struct {
			Min float64 `yaml:"min"`
			Max float64 `yaml:"max"`
		} `yaml:"summer"`
		Winter struct {
			Min float64 `yaml:"min"`
			Max float64 `yaml:"max"`
		} `yaml:"winter"`
	} `yaml:"pressure"`
}

var topics = brokerconfiguration.GetAlertManagerTopics()

var (
	BROKER      = brokerconfiguration.GetBrokerAddress()
	TOPIC       = topics[0]
	ALERT_TOPIC = topics[1]
)

func getThresholds() (Thresholds, error) {
	yamlFile, err := os.Open("config/threshold_config.yml")
	if err != nil {
		log.Fatal("Error opening threshold configuration file:", err)
		return Thresholds{}, err
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)

	var thresholds Thresholds

	err = yaml.Unmarshal(byteValue, &thresholds)
	if err != nil {
		log.Fatal("Error unmarshalling YAML:", err)
		return Thresholds{}, err
	}

	return thresholds, nil
}

func getAlertTopicFromMessageTopic(topic string) string {
	return ALERT_TOPIC + brokerutils.GetAirportCodeFromTopic(topic)
}

func getSeasonFromTimestamp(timestamp string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
		return ""
	}

	month := t.Month()

	if month >= 3 && month <= 9 {
		return "summer"
	} else {
		return "winter"
	}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	thresholds, err := getThresholds()
	if err != nil {
		log.Println("Error getting thresholds:", err)
		return
	}

	payload := string(message.Payload())
	data := strings.Split(payload, " ")

	submittedTimestamp := data[0] + " " + data[1]
	sensor := data[2]
	value, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		log.Println("Failed to convert value to integer", err)
		return
	}

	switch sensor {
	case "temperature":
		if value < thresholds.Temp.Min || value > thresholds.Temp.Max {
			alertMessage := fmt.Sprintf("Alert: Temperature (%f) exceeded threshold (%f-%f)", value, thresholds.Temp.Min, thresholds.Temp.Max)
			token := client.Publish(getAlertTopicFromMessageTopic(message.Topic()), 1, false, alertMessage)
			token.Wait()
		}
	case "pressure":
		season := getSeasonFromTimestamp(submittedTimestamp)
		var minThreshold, maxThreshold float64
		if season == "summer" {
			minThreshold, maxThreshold = thresholds.Pressure.Summer.Min, thresholds.Pressure.Summer.Max
		} else if season == "winter" {
			minThreshold, maxThreshold = thresholds.Pressure.Winter.Min, thresholds.Pressure.Winter.Max
		} else {
			log.Println("Error getting season from timestamp")
			return
		}

		if value < minThreshold || value > maxThreshold {
			alertMessage := fmt.Sprintf("Alert: Pressure (%f) exceeded threshold (%f-%f)", value, minThreshold, maxThreshold)
			token := client.Publish(getAlertTopicFromMessageTopic(message.Topic()), 1, false, alertMessage)
			token.Wait()
		}
	case "wind":
		if value > thresholds.Wind.Speed {
			alertMessage := fmt.Sprintf("Alert: Wind (%f) exceeded threshold (%f)", value, thresholds.Wind.Speed)
			token := client.Publish(getAlertTopicFromMessageTopic(message.Topic()), 1, false, alertMessage)
			token.Wait()
		}
	default:
		log.Fatalf("Unknown sensor: %s\n", sensor)
	}

}

func main() {
	client, err := mqttconnect.NewClient(BROKER, "alert_manager", onMessageReceived)
	if err != nil {
		log.Fatal("Error creating MQTT client:", err)
		return
	}

	err = client.Subscribe(TOPIC, 1, nil)
	if err != nil {
		log.Println("Error subscribing to topic:", err)
		return
	}

	mqttconnect.WaitForSignal()
}

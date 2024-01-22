package brokerconfiguration

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Config struct {
	BrokerAddress string `yaml:"brokerAddress"`
	Topics        struct {
		AlertManager struct {
			Subscribe string `yaml:"subscribe"`
			Publish   string `yaml:"publish"`
		} `yaml:"alertManager"`
	} `yaml:"topics"`
	InfluxDB struct {
		Bucket    string `yaml:"bucket"`
		Org       string `yaml:"org"`
		Url       string `yaml:"url"`
		Subscribe string `yaml:"subscribe"`
	} `yaml:"influxdb"`
	FileRecorder struct {
		Subscribe     string `yaml:"subscribe"`
		RecordingPath string `yaml:"recordingPath"`
	} `yaml:"fileRecorder"`
}

func getAppConfig() (Config, error) {
	yamlFile, err := os.Open("config/app_config.yml")
	if err != nil {
		log.Fatal("Error reading app config file:", err)
		return Config{}, err
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)

	var config Config
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal("Error unmarshalling YAML:", err)
		return Config{}, err
	}
	return config, nil
}

func GetBrokerAddress() string {
	config, err := getAppConfig()
	if err != nil {
		log.Fatalf("Error getting app config: %v", err)
		return ""
	}

	brokerAddress := config.BrokerAddress

	return brokerAddress
}

func GetAlertManagerTopics() []string {
	config, err := getAppConfig()
	if err != nil {
		log.Fatalf("Error getting app config: %v", err)
		return []string{}
	}

	alertManagerTopicSubscribe := config.Topics.AlertManager.Subscribe
	alertManagerTopicPublish := config.Topics.AlertManager.Publish

	alertManagerTopics := []string{alertManagerTopicSubscribe, alertManagerTopicPublish}

	return alertManagerTopics
}

func GetInfluxdbSettings() []string {
	config, err := getAppConfig()
	if err != nil {
		log.Fatalf("Error getting app config: %v", err)
		return []string{}
	}

	influxdbBucket := config.InfluxDB.Bucket
	influxdbOrg := config.InfluxDB.Org
	influxdbUrl := config.InfluxDB.Url
	influxdbSubscribe := config.InfluxDB.Subscribe

	databaseRecorderConfig := []string{influxdbBucket, influxdbOrg, influxdbUrl, influxdbSubscribe}

	return databaseRecorderConfig
}

func GetFileRecorderSettings() []string {
	config, err := getAppConfig()
	if err != nil {
		log.Fatalf("Error getting app config: %v", err)
		return []string{}
	}

	fileRecorderTopic := config.FileRecorder.Subscribe
	fileRecorderFolder := config.FileRecorder.RecordingPath

	fileRecorderConfig := []string{fileRecorderTopic, fileRecorderFolder}

	return fileRecorderConfig
}

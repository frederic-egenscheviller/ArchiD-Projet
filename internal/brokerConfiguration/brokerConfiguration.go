package brokerConfiguration

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BrokerAddress string `yaml:"brokerAddress"`
	Topics        struct {
		AlertManager struct {
			Subscribe string `yaml:"subscribe"`
			Publish   string `yaml:"publish"`
		} `yaml:"alertManager"`
	} `yaml:"topics"`
}

func getAppConfig() (Config, error) {
	yamlFile, err := os.Open("config/app_config.yml")
	if err != nil {
		fmt.Println("Error reading app config file:", err)
		return Config{}, err
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)

	var config Config
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return Config{}, err
	}
	return config, nil
}

func GetBrokerAddress() string {
	config, err := getAppConfig()
	if err != nil {
		fmt.Printf("Error getting app config: %v", err)
		return ""
	}

	brokerAddress := config.BrokerAddress

	return brokerAddress
}

func GetAlertManagerTopics() ([]string) {
	config, err := getAppConfig()
	if err != nil {
		fmt.Printf("Error getting app config: %v", err)
		return []string{}
	}

	alertManagerTopicSubscribe := config.Topics.AlertManager.Subscribe
	alertManagerTopicPublish := config.Topics.AlertManager.Publish

	alertManagerTopics := []string{alertManagerTopicSubscribe, alertManagerTopicPublish}

	return alertManagerTopics
}
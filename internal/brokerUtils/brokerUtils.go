package brokerutils

import (
	"strings"
)

func GetAirportCodeFromTopic(topic string) string {
	return strings.Split(topic, "/")[1]
}

package meteofranceAPI

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type APIResponse struct {
	Lat           float64   `json:"lat"`
	Lon           float64   `json:"lon"`
	GeoIDInsee    string    `json:"geo_id_insee"`
	ReferenceTime time.Time `json:"reference_time"`
	InsertTime    time.Time `json:"insert_time"`
	ValidityTime  time.Time `json:"validity_time"`
	T             float64   `json:"t"`
	Ff            float64   `json:"ff"`
	Pres          int       `json:"pres"`
}

type SensorData struct {
	SensorID         int
	AirportID        string
	Measurement      string
	MeasurementValue float64
	MeasurementTime  time.Time
}

type SensorInfo struct {
	TransmissionFrequency time.Duration `yaml:"transmissionFrequency"`
	ClientID              string        `yaml:"clientID"`
	QoS                   byte          `yaml:"qos"`
	AirportIATA           string        `yaml:"airportIATA"`
	GeoIDInsee            string        `yaml:"geoIDInsee"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func FetchSensorDataFromAPI(sensorInfo SensorInfo) (SensorData, error) {
	apiKey := os.Getenv("METEO_FRANCE_API_KEY")

	apiURL := fmt.Sprintf("https://public-api.meteofrance.fr/public/DPObs/v1/station/infrahoraire-6m?id_station=%s&format=json", sensorInfo.GeoIDInsee)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return SensorData{}, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Add("apikey", apiKey)
	req.Header.Add("accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SensorData{}, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SensorData{}, fmt.Errorf("error reading response body: %v", err)
	}

	var apiResponses []APIResponse
	err = json.Unmarshal(body, &apiResponses)
	if err != nil {
		return SensorData{}, fmt.Errorf("error decoding JSON response: %v", err)
	}

	apiResponse := apiResponses[0]

	referenceTimeUTC1 := apiResponse.ReferenceTime.In(time.FixedZone("UTC+1", 60*60))

	var sensorData SensorData

	switch true {
	case strings.HasPrefix(sensorInfo.ClientID, "pressure_sensor"):
		sensorData = SensorData{
			SensorID:         1,
			AirportID:        sensorInfo.AirportIATA,
			Measurement:      "pressure",
			MeasurementValue: float64(apiResponse.Pres) / 100,
			MeasurementTime:  referenceTimeUTC1,
		}
	case strings.HasPrefix(sensorInfo.ClientID, "temperature_sensor"):
		sensorData = SensorData{
			SensorID:         1,
			AirportID:        sensorInfo.AirportIATA,
			Measurement:      "temperature",
			MeasurementValue: apiResponse.T,
			MeasurementTime:  referenceTimeUTC1,
		}
	case strings.HasPrefix(sensorInfo.ClientID, "wind_sensor"):
		sensorData = SensorData{
			SensorID:         1,
			AirportID:        sensorInfo.AirportIATA,
			Measurement:      "wind",
			MeasurementValue: apiResponse.Ff,
			MeasurementTime:  referenceTimeUTC1,
		}
	default:
		return SensorData{}, fmt.Errorf("unknown measurement type: %s", sensorInfo.ClientID)
	}
	return sensorData, nil
}

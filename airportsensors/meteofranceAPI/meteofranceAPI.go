package meteofranceAPI

import (
	"ArchiD-Projet/internal/sensors"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"time"
)

type APIResponse struct {
	Lat           float64     `json:"lat"`
	Lon           float64     `json:"lon"`
	GeoIDInsee    string      `json:"geo_id_insee"`
	ReferenceTime time.Time   `json:"reference_time"`
	InsertTime    time.Time   `json:"insert_time"`
	ValidityTime  time.Time   `json:"validity_time"`
	T             float64     `json:"t"`
	U             int         `json:"u"`
	Dd            int         `json:"dd"`
	Ff            float64     `json:"ff"`
	Dxi10         int         `json:"dxi10"`
	Fxi10         float64     `json:"fxi10"`
	RrPer         int         `json:"rr_per"`
	T10           float64     `json:"t_10"`
	T20           float64     `json:"t_20"`
	T50           float64     `json:"t_50"`
	T100          float64     `json:"t_100"`
	Vv            int         `json:"vv"`
	EtatSol       interface{} `json:"etat_sol"`
	Sss           int         `json:"sss"`
	N             int         `json:"n"`
	Insolh        int         `json:"insolh"`
	RayGlo01      int         `json:"ray_glo01"`
	Pres          int         `json:"pres"`
	Pmer          interface{} `json:"pmer"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func FetchSensorDataFromAPI(measurementType string) (sensors.SensorData, error) {
	apiKey := os.Getenv("METEO_FRANCE_API_KEY")

	apiURL := "https://public-api.meteofrance.fr/public/DPObs/v1/station/infrahoraire-6m?id_station=13054001&format=json"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return sensors.SensorData{}, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Add("apikey", apiKey)
	req.Header.Add("accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return sensors.SensorData{}, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sensors.SensorData{}, fmt.Errorf("error reading response body: %v", err)
	}

	var apiResponses []APIResponse
	err = json.Unmarshal(body, &apiResponses)
	if err != nil {
		return sensors.SensorData{}, fmt.Errorf("error decoding JSON response: %v", err)
	}

	apiResponse := apiResponses[0]

	referenceTimeUTC1 := apiResponse.ReferenceTime.In(time.FixedZone("UTC+1", 60*60))

	var sensorData sensors.SensorData

	switch measurementType {
	case "pressure_sensor":
		sensorData = sensors.SensorData{
			SensorID:         1,
			AirportID:        "MRS",
			Measurement:      "pressure",
			MeasurementValue: float64(apiResponse.Pres) / 100,
			MeasurementTime:  referenceTimeUTC1,
		}
	case "temperature_sensor":
		sensorData = sensors.SensorData{
			SensorID:         1,
			AirportID:        "MRS",
			Measurement:      "temperature",
			MeasurementValue: apiResponse.T,
			MeasurementTime:  referenceTimeUTC1,
		}
	case "wind_sensor":
		sensorData = sensors.SensorData{
			SensorID:         1,
			AirportID:        "MRS",
			Measurement:      "wind",
			MeasurementValue: apiResponse.Ff,
			MeasurementTime:  referenceTimeUTC1,
		}
	default:
		return sensors.SensorData{}, fmt.Errorf("unknown measurement type: %s", measurementType)
	}

	return sensorData, nil
}

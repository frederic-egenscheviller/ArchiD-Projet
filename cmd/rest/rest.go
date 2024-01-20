package main

import (
	brokerconfiguration "ArchiD-Projet/internal/brokerConfiguration"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"time"
)

var influxDBAPIKey string
var influxDBURL string
var influxDBBucket string
var influxDBOrg string
var influxDBClient influxdb2.Client
var loc = time.Local

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := brokerconfiguration.GetInfluxdbSettings()

	influxDBAPIKey = os.Getenv("INFLUX_DB_API_KEY")
	influxDBBucket = config[0]
	influxDBOrg = config[1]
	influxDBURL = config[2]

	loc, err = time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal("Error loading timezone")
		return
	}

	if influxDBAPIKey == "" || influxDBURL == "" || influxDBBucket == "" {
		log.Fatal("Incomplete InfluxDB configuration in app_config.yml")
	}

	influxDBClient = influxdb2.NewClientWithOptions(influxDBURL, influxDBAPIKey, influxdb2.DefaultOptions())
}

type data struct {
	AirportIATA string  `json:"airport"`
	Datetime    string  `json:"_time"`
	Type        string  `json:"_measurement"`
	Value       float64 `json:"_value"`
}

type dataAverage struct {
	AirportIATA string  `json:"airport"`
	Measurement string  `json:"measurement"`
	Value       float64 `json:"value"`
}

type sensor struct {
	AirportIATA string `json:"airport"`
	Measurement string `json:"type"`
}

type airport struct {
	AirportIATA string `json:"airport"`
}

func main() {
	defer influxDBClient.Close()

	router := gin.Default()
	router.GET("/airports", getAllAirport)
	router.GET("/airports/data/", getAllAirportData)
	router.GET("/airport/:id/data/", getAirportDataById)
	router.GET("/airport/:id/sensors", getSensorsByAirportId)
	router.GET("/airport/:id/data/range/:start/:end/:type", getAirportDataByDateRangesAndType)
	router.GET("/airport/:id/average/:date", getAirportDataAverageByDate)
	router.GET("/airport/:id/average/:date/:type", getAirportDataAverageByDateAndType)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Error starting Gin router:", err)
	}
}

func getAllAirport(c *gin.Context) {
	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> group(columns: ["airport"]) |> distinct(column: "airport")`, influxDBBucket)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []airport
	for result.Next() {
		value := result.Record().ValueByKey("airport")
		if value != nil {
			airportID, ok := value.(string)
			if ok {
				ret = append(ret, airport{AirportIATA: airportID})
			}
		}
	}
	c.IndentedJSON(http.StatusOK, ret)
}

func getSensorsByAirportId(c *gin.Context) {
	id := c.Param("id")

	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> filter(fn: (r) => r["airport"] == "%s") |> distinct(column: "_measurement")`, influxDBBucket, id)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sensor data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []sensor
	for result.Next() {
		value := result.Record().ValueByKey("_measurement")
		if value != nil {
			sensorType, ok := value.(string)
			if ok {
				ret = append(ret, sensor{AirportIATA: id, Measurement: sensorType})
			}
		}
	}

	c.IndentedJSON(http.StatusOK, ret)
}

func getAllAirportData(c *gin.Context) {
	query := fmt.Sprintf(`
        from(bucket:"%s") 
        |> range(start: 1970-01-01T00:00:00Z)`,
		influxDBBucket)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []data
	for result.Next() {
		airportIATA := result.Record().ValueByKey("airport")
		datetime := result.Record().ValueByKey("_time")
		sensorType := result.Record().ValueByKey("_measurement")
		value := result.Record().ValueByKey("_value")
		datetimeUTC1 := datetime.(time.Time).In(loc).Format("2006-01-02 15:04:05")

		if airportIATA != nil && datetime != nil && sensorType != nil && value != nil {
			ret = append(ret, data{
				AirportIATA: airportIATA.(string),
				Datetime:    datetimeUTC1,
				Type:        sensorType.(string),
				Value:       value.(float64),
			})
		}
	}

	c.IndentedJSON(http.StatusOK, ret)
}

func getAirportDataById(c *gin.Context) {
	id := c.Param("id")

	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> filter(fn: (r) => r["airport"] == "%s")`, influxDBBucket, id)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []data
	for result.Next() {
		airportIATA := result.Record().ValueByKey("airport")
		datetime := result.Record().ValueByKey("_time")
		sensorType := result.Record().ValueByKey("_measurement")
		value := result.Record().ValueByKey("_value")
		datetimeUTC1 := datetime.(time.Time).In(loc).Format("2006-01-02 15:04:05")

		if airportIATA != nil && datetime != nil && sensorType != nil && value != nil {
			ret = append(ret, data{
				AirportIATA: airportIATA.(string),
				Datetime:    datetimeUTC1,
				Type:        sensorType.(string),
				Value:       value.(float64),
			})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified airport ID"})
}

func getAirportDataByDateRangesAndType(c *gin.Context) {
	id := c.Param("id")
	dataType := c.Param("type")
	start := c.Param("start")
	end := c.Param("end")

	query := fmt.Sprintf(`
        from(bucket: "%s")
  			|> range(start: %s, stop: %s)
  			|> filter(fn: (r) => r["airport"] == "%s" and r["_measurement"] == "%s")`,
		influxDBBucket, start, end, id, dataType)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []data
	for result.Next() {
		airportIATA := result.Record().ValueByKey("airport")
		datetime := result.Record().ValueByKey("_time")
		sensorType := result.Record().ValueByKey("_measurement")
		value := result.Record().ValueByKey("_value")
		datetimeUTC1 := datetime.(time.Time).In(loc).Format("2006-01-02 15:04:05")

		if airportIATA != nil && datetime != nil && sensorType != nil && value != nil {
			ret = append(ret, data{
				AirportIATA: airportIATA.(string),
				Datetime:    datetimeUTC1,
				Type:        sensorType.(string),
				Value:       value.(float64),
			})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

func getAirportDataAverageByDate(c *gin.Context) {
	id := c.Param("id")
	startDate := c.Param("date")

	startDateFomatted, err := time.Parse("2006-01-02", startDate)

	endDate := startDateFomatted.AddDate(0, 0, 1).Format("2006-01-02")

	startDate = startDate + "T00:00:00Z"
	endDate = endDate + "T00:00:00Z"

	var ret []dataAverage

	query := fmt.Sprintf(`
		from(bucket: "%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r["airport"] == "%s")
		  |> group(columns: ["_measurement"])
		  |> aggregateWindow(every: 1d, fn: mean, createEmpty: false)`,
		influxDBBucket, startDate, endDate, id)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	for result.Next() {
		measurement := result.Record().ValueByKey("_measurement")
		value := result.Record().ValueByKey("_value")
		if value != nil {
			ret = append(ret, dataAverage{AirportIATA: id, Measurement: measurement.(string), Value: value.(float64)})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

func getAirportDataAverageByDateAndType(c *gin.Context) {
	id := c.Param("id")
	startDate := c.Param("date")
	measurement := c.Param("type")

	startDateFomatted, err := time.Parse("2006-01-02", startDate)

	endDate := startDateFomatted.AddDate(0, 0, 1).Format("2006-01-02")

	startDate = startDate + "T00:00:00Z"
	endDate = endDate + "T00:00:00Z"

	var ret []dataAverage

	query := fmt.Sprintf(`
		from(bucket: "%s")
		  |> range(start: %s, stop: %s)
		  |> filter(fn: (r) => r["airport"] == "%s" and r["_measurement"] == "%s")
		  |> aggregateWindow(every: 1d, fn: mean, createEmpty: false)`,
		influxDBBucket, startDate, endDate, id, measurement)

	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	for result.Next() {
		measurement := result.Record().ValueByKey("_measurement")
		value := result.Record().ValueByKey("_value")
		if value != nil {
			ret = append(ret, dataAverage{AirportIATA: id, Measurement: measurement.(string), Value: value.(float64)})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

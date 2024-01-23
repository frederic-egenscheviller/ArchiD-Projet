package main

// @title ArchiD-Projet API
// @description This is the API for the ArchiD-Projet project
// @version 1
// @host localhost:8080
// @BasePath /

import (
	"ArchiD-Projet/cmd/rest/docs"
	brokerconfiguration "ArchiD-Projet/internal/brokerConfiguration"
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var influxDBAPIKey string
var influxDBURL string
var influxDBBucket string
var influxDBOrg string
var influxDBClient influxdb2.Client
var loc = time.Local

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(currentFile)

	// Construct the absolute path to the .env file
	envFilePath := filepath.Join(projectRoot+"/../../", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	config := brokerconfiguration.GetInfluxdbSettings()

	influxDBAPIKey = os.Getenv("INFLUX_DB_API_KEY")
	influxDBBucket = config[0]
	influxDBOrg = config[1]
	influxDBURL = config[2]

	loc, err = time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Println("Error loading timezone", err)
		return
	}

	if influxDBAPIKey == "" || influxDBURL == "" || influxDBBucket == "" {
		log.Fatal("Incomplete InfluxDB configuration in app_config.yml")
	}

	influxDBClient = influxdb2.NewClientWithOptions(influxDBURL, influxDBAPIKey, influxdb2.DefaultOptions())
}

type data struct {
	AirportIATA string  `json:"airport"`
	Datetime    string  `json:"time"`
	Type        string  `json:"measurement"`
	Value       float64 `json:"value"`
}

type dataAverage struct {
	AirportIATA string  `json:"airport"`
	Measurement string  `json:"measurement"`
	Value       float64 `json:"value"`
}

type sensor struct {
	AirportIATA string `json:"airport"`
	Measurement string `json:"measurement"`
}

type airport struct {
	AirportIATA string `json:"airport"`
}

func main() {
	defer influxDBClient.Close()
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	router.GET("/airports", getAllAirports)
	router.GET("/airports/data/", getAllAirportsData)
	router.GET("/airport/:iata/data/", getAirportDataByIATA)
	router.GET("/airport/:iata/sensors", getSensorsByAirportIATA)
	router.GET("/airport/:iata/data/range/:start/:end/:measurement", getAirportDataByDateRangesAndType)
	router.GET("/airport/:iata/average/:date", getAirportDataAverageByDate)
	router.GET("/airport/:iata/average/:date/:measurement", getAirportDataAverageByDateAndType)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Error starting Gin router:", err)
	}
}

// @BasePath /
// @Summary Get all airports
// @Description Get all airports
// @Accept json
// @Produce json
// @Success 200 {array} airport
// @Router /airports [get]
func getAllAirports(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
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
			airportIATA, ok := value.(string)
			if ok {
				ret = append(ret, airport{AirportIATA: airportIATA})
			}
		}
	}
	c.IndentedJSON(http.StatusOK, ret)
}

// @BasePath /
// @Summary Get all sensors for an airport
// @Description Get all sensors for an airport
// @Accept json
// @Produce json
// @Param iata path string true "Airport IATA code"
// @Success 200 {array} sensor
// @Router /airport/{iata}/sensors [get]
func getSensorsByAirportIATA(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	airportIATA := c.Param("iata")

	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> filter(fn: (r) => r["airport"] == "%s") |> distinct(column: "_measurement")`, influxDBBucket, airportIATA)

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
				ret = append(ret, sensor{AirportIATA: airportIATA, Measurement: sensorType})
			}
		}
	}
	c.IndentedJSON(http.StatusOK, ret)
}

// @BasePath /
// @Summary Get all data for all airports
// @Description Get all data for all airports
// @Accept json
// @Produce json
// @Success 200 {array} data
// @Router /airports/data [get]
func getAllAirportsData(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
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

// @BasePath /
// @Summary Get all data for an airport
// @Description Get all data for an airport
// @Accept json
// @Produce json
// @Param iata path string true "Airport IATA code"
// @Success 200 {array} data
// @Router /airport/{iata}/data [get]
func getAirportDataByIATA(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	airportIATA := c.Param("iata")

	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> filter(fn: (r) => r["airport"] == "%s")`, influxDBBucket, airportIATA)

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

// @BasePath /
// @Summary Get all data for an airport between two dates
// @Description Get all data for an airport between two dates
// @Accept json
// @Produce json
// @Param iata path string true "Airport IATA code"
// @Param start path string true "Start date"
// @Param end path string true "End date"
// @Param measurement path string true "Measurement type"
// @Success 200 {array} data
// @Router /airport/{iata}/data/range/{start}/{end}/{measurement} [get]
func getAirportDataByDateRangesAndType(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")

	airportIATA := c.Param("iata")
	dataType := c.Param("measurement")
	start := c.Param("start")
	end := c.Param("end")

	query := fmt.Sprintf(`
        from(bucket: "%s")
  			|> range(start: %s, stop: %s)
  			|> filter(fn: (r) => r["airport"] == "%s" and r["_measurement"] == "%s")`,
		influxDBBucket, start, end, airportIATA, dataType)

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

// @BasePath /
// @Summary Get average data for an airport for a date
// @Description Get average data for an airport for a date
// @Accept json
// @Produce json
// @Param iata path string true "Airport IATA code"
// @Param date path string true "Date"
// @Success 200 {array} dataAverage
// @Router /airport/{iata}/average/{date} [get]
func getAirportDataAverageByDate(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")

	airportIATA := c.Param("iata")
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
		influxDBBucket, startDate, endDate, airportIATA)

	fmt.Println(query)

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
			ret = append(ret, dataAverage{AirportIATA: airportIATA, Measurement: measurement.(string), Value: value.(float64)})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

// @BasePath /
// @Summary Get average data for an airport for a date and a measurement type
// @Description Get average data for an airport for a date and a measurement type
// @Accept json
// @Produce json
// @Param iata path string true "Airport IATA code"
// @Param date path string true "Date"
// @Param measurement path string true "Measurement type"
// @Success 200 {array} dataAverage
// @Router /airport/{iata}/average/{date}/{measurement} [get]
func getAirportDataAverageByDateAndType(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")

	airportIATA := c.Param("iata")
	startDate := c.Param("date")
	dataType := c.Param("measurement")

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
		influxDBBucket, startDate, endDate, airportIATA, dataType)

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
			ret = append(ret, dataAverage{AirportIATA: airportIATA, Measurement: measurement.(string), Value: value.(float64)})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

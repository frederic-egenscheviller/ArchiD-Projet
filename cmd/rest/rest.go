package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
)

var influxDBAPIKey string
var influxDBURL string
var influxDBBucket string
var influxDBOrg string
var influxDBClient influxdb2.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	influxDBAPIKey = os.Getenv("INFLUX_DB_API_KEY")
	influxDBURL = os.Getenv("INFLUX_DB_URL")
	influxDBBucket = os.Getenv("INFLUX_DB_BUCKET")
	influxDBOrg = os.Getenv("INFLUX_DB_ORG")

	if influxDBAPIKey == "" || influxDBURL == "" || influxDBBucket == "" {
		log.Fatal("Incomplete InfluxDB configuration in .env")
	}

	// Create the InfluxDB client once
	influxDBClient = influxdb2.NewClientWithOptions(influxDBURL, influxDBAPIKey, influxdb2.DefaultOptions())
}

type data struct {
	AirportId string  `json:"id"`
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
}

type dataAverage struct {
	AirportId string  `json:"id"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
}

type sensor struct {
	AirportId string `json:"id"`
	Type      string `json:"type"`
}

type airport struct {
	AirportId string `json:"id"`
}

// Je sais pas si on en a encore besoin
var dataArr []data

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

func setDataArray(slice []data) {
	dataArr = slice
}

func getAllAirport(c *gin.Context) {
	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> distinct(column: "AirportId")`, influxDBBucket)

	// Query data from InfluxDB using the global client
	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []airport
	for result.Next() {
		// Extract values from the result
		value := result.Record().ValueByKey("AirportId")
		if value != nil {
			airportID, ok := value.(string)
			if ok {
				ret = append(ret, airport{AirportId: airportID})
			}
		}
	}

	// Check if any airports were found
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No airport found"})
}

func getSensorsByAirportId(c *gin.Context) {
	id := c.Param("id")

	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> filter(fn: (r) => r._measurement == "your_measurement" and r.AirportId == "%s") |> distinct(column: "Type")`, influxDBBucket, id)

	// Query data from InfluxDB using the global client
	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sensor data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []sensor
	for result.Next() {
		// Extract values from the result
		value := result.Record().ValueByKey("Type")
		if value != nil {
			sensorType, ok := value.(string)
			if ok {
				ret = append(ret, sensor{AirportId: id, Type: sensorType})
			}
		}
	}

	// Check if any sensors were found
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No sensor found for the specified airport ID"})
}

func getAllAirportData(c *gin.Context) {
	query := fmt.Sprintf(`
        from(bucket:"%s") 
        |> range(start: 1970-01-01T00:00:00Z) 
        |> group(columns: ["AirportId"])
        |> last()`,
		influxDBBucket)

	// Query data from InfluxDB using the global client
	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []data
	for result.Next() {
		// Extract values from the result
		airportID := result.Record().ValueByKey("AirportId")
		date := result.Record().ValueByKey("date")
		time := result.Record().ValueByKey("time")
		sensorType := result.Record().ValueByKey("Type")
		value := result.Record().ValueByKey("value")

		if airportID != nil && date != nil && time != nil && sensorType != nil && value != nil {
			ret = append(ret, data{
				AirportId: airportID.(string),
				Date:      date.(string),
				Time:      time.(string),
				Type:      sensorType.(string),
				Value:     value.(float64),
			})
		}
	}

	// Check if any data were found
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for any airport"})
}

func getAirportDataById(c *gin.Context) {
	id := c.Param("id")

	query := fmt.Sprintf(`from(bucket:"%s") |> range(start: 1970-01-01T00:00:00Z) |> filter(fn: (r) => r.AirportId == "%s")`, influxDBBucket, id)

	// Query data from InfluxDB using the global client
	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []data
	for result.Next() {
		// Extract values from the result
		airportID := result.Record().ValueByKey("AirportId")
		date := result.Record().ValueByKey("date")
		time := result.Record().ValueByKey("time")
		sensorType := result.Record().ValueByKey("Type")
		value := result.Record().ValueByKey("value")

		if airportID != nil && date != nil && time != nil && sensorType != nil && value != nil {
			ret = append(ret, data{
				AirportId: airportID.(string),
				Date:      date.(string),
				Time:      time.(string),
				Type:      sensorType.(string),
				Value:     value.(float64),
			})
		}
	}

	// Check if any data were found
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
        from(bucket:"%s") 
        |> range(start: %s, stop: %s) 
        |> filter(fn: (r) => r.AirportId == "%s" and r.Type == "%s")`,
		influxDBBucket, start, end, id, dataType)

	// Query data from InfluxDB using the global client
	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []data
	for result.Next() {
		// Extract values from the result
		airportID := result.Record().ValueByKey("AirportId")
		date := result.Record().ValueByKey("date")
		time := result.Record().ValueByKey("time")
		sensorType := result.Record().ValueByKey("Type")
		value := result.Record().ValueByKey("value")

		if airportID != nil && date != nil && time != nil && sensorType != nil && value != nil {
			ret = append(ret, data{
				AirportId: airportID.(string),
				Date:      date.(string),
				Time:      time.(string),
				Type:      sensorType.(string),
				Value:     value.(float64),
			})
		}
	}

	// Check if any data were found
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

func getAirportDataAverageByDate(c *gin.Context) {
	id := c.Param("id")
	date := c.Param("date")

	var dataTypes = []string{"temperature", "wind", "pressure"}
	var ret []dataAverage

	for _, dataType := range dataTypes {
		query := fmt.Sprintf(`
            from(bucket:"%s") 
            |> range(start: %s, stop: %s) 
            |> filter(fn: (r) => r.AirportId == "%s" and r.Type == "%s" and r.date == "%s")
            |> mean(column: "value")`,
			influxDBBucket, date, date, id, dataType, date)

		// Query data from InfluxDB using the global client
		result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
			return
		}
		defer result.Close()

		for result.Next() {
			// Extract values from the result
			value := result.Record().Value()
			if value != nil {
				ret = append(ret, dataAverage{AirportId: id, Type: dataType, Value: value.(float64)})
			}
		}
	}

	// Check if any data were found
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

func getAirportDataAverageByDateAndType(c *gin.Context) {
	id := c.Param("id")
	date := c.Param("date")
	dataType := c.Param("type")

	query := fmt.Sprintf(`
        from(bucket:"%s") 
        |> range(start: %s, stop: %s) 
        |> filter(fn: (r) => r.AirportId == "%s" and r.Type == "%s" and r.date == "%s")
        |> mean(column: "value")`,
		influxDBBucket, date, date, id, dataType, date)

	// Query data from InfluxDB using the global client
	result, err := influxDBClient.QueryAPI(influxDBOrg).Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from InfluxDB"})
		return
	}
	defer result.Close()

	var ret []dataAverage
	for result.Next() {
		// Extract values from the result
		value := result.Record().Value()
		if value != nil {
			ret = append(ret, dataAverage{AirportId: id, Type: dataType, Value: value.(float64)})
		}
	}

	// Check if any data were found
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found for the specified parameters"})
}

// Je sais pas si on en a encore besoin
func getValueListByType(d []data, dataType string) []data {
	var ret []data
	for _, a := range d {
		if a.Type == dataType {
			ret = append(ret, a)
		}
	}
	return ret
}

// Je sais pas si on en a encore besoin
func getAverageValue(data []data) float64 {
	var sum float64
	for _, a := range data {
		sum += a.Value
	}
	return sum / float64(len(data))
}

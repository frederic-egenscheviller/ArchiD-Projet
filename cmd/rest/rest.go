package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

var dataArr = []data{
	{AirportId: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "temperature", Value: 20.0},
	{AirportId: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "wind", Value: 3.100000},
	{AirportId: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "pressure", Value: 1024.0},
	{AirportId: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "temperature", Value: 24.3},
	{AirportId: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "wind", Value: 2.0},
	{AirportId: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "pressure", Value: 1020.0},
}

func main() {
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
		return
	}
}

func containsElement(slice []airport, element string) bool {
	for _, value := range slice {
		if value.AirportId == element {
			return true
		}
	}
	return false
}

func getAllAirport(c *gin.Context) {
	var ret []airport
	for _, a := range dataArr {
		if !containsElement(ret, a.AirportId) {
			ret = append(ret, airport{AirportId: a.AirportId})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No airport found"})
}

func getSensorsByAirportId(c *gin.Context) {
	id := c.Param("id")
	var ret []sensor

	for _, a := range dataArr {
		if a.AirportId == id {
			ret = append(ret, sensor{AirportId: a.AirportId, Type: a.Type})
		}
	}

	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No sensor found"})
}

func getAllAirportData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dataArr)

}

func getAirportDataById(c *gin.Context) {
	id := c.Param("id")
	var ret []data

	for _, a := range dataArr {
		if a.AirportId == id {
			ret = append(ret, a)
		}
	}
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "airport not found"})
}

func getAirportDataByDateRangesAndType(c *gin.Context) {
	id := c.Param("id")
	dataType := c.Param("type")
	start := c.Param("start")
	end := c.Param("end")

	var ret []data

	for _, a := range dataArr {
		if a.AirportId == id && a.Type == dataType && a.Date >= start && a.Date <= end {
			ret = append(ret, a)
		}
	}
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found"})
}

func getAirportDataAverageByDate(c *gin.Context) {
	id := c.Param("id")
	date := c.Param("date")

	var dataList []data
	var ret []dataAverage
	var dataTypes = []string{"temperature", "wind", "pressure"}

	for _, a := range dataArr {
		if a.AirportId == id && a.Date == date {
			dataList = append(dataList, a)
		}
	}
	if len(dataList) != 0 {
		for _, a := range dataTypes {
			ret = append(ret, dataAverage{AirportId: id, Type: a, Value: getAverageValue(getValueListByType(dataList, a))})
		}
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found"})
}

func getAirportDataAverageByDateAndType(c *gin.Context) {
	id := c.Param("id")
	date := c.Param("date")
	dataType := c.Param("type")

	var ret []data

	for _, a := range dataArr {
		if a.AirportId == id && a.Date == date && a.Type == dataType {
			ret = append(ret, a)
		}
	}
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"AirportId": id, "Date": date, dataType + "Average": getAverageValue(ret)})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No data found"})
}

func getValueListByType(d []data, dataType string) []data {
	var ret []data
	for _, a := range d {
		if a.Type == dataType {
			ret = append(ret, a)
		}
	}
	return ret
}

func getAverageValue(data []data) float64 {
	var sum float64
	for _, a := range data {
		sum += a.Value
	}
	return sum / float64(len(data))
}

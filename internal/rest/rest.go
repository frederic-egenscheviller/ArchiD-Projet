package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type data struct {
	AirportID string  `json:"id"`
	Date      string  `json:"date"`
	Time      string  `json:"time"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
}

type dataAverage struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

var dataArr = []data{
	{AirportID: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "temperature", Value: 20.0},
	{AirportID: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "wind", Value: 3.100000},
	{AirportID: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "pressure", Value: 1024.0},
	{AirportID: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "temperature", Value: 24.3},
	{AirportID: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "wind", Value: 2.0},
	{AirportID: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "pressure", Value: 1020.0},
}

func main() {
	router := gin.Default()
	router.GET("/airport-data", getAllAirportData)
	router.GET("/airport-data/:id", getAirportDataById)
	router.GET("/airport-data/:id/:type/:start/:end", getAirportDataByTypeAndDateRanges)
	router.GET("/airport-data/average/:id/:date", getAirportDataAverageByDate)
	router.GET("/airport-data/average/:id/:date/:type", getAirportDataAverageByDateAndType)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}

func getAllAirportData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dataArr)
}

func getAirportDataById(c *gin.Context) {
	id := c.Param("id")
	var ret []data

	for _, a := range dataArr {
		if a.AirportID == id {
			ret = append(ret, a)
		}
	}
	if len(ret) != 0 {
		c.IndentedJSON(http.StatusOK, ret)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "airport not found"})
}

func getAirportDataByTypeAndDateRanges(c *gin.Context) {
	id := c.Param("id")
	dataType := c.Param("type")
	start := c.Param("start")
	end := c.Param("end")

	var ret []data

	for _, a := range dataArr {
		if a.AirportID == id && a.Type == dataType && a.Date >= start && a.Date <= end {
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
		if a.AirportID == id && a.Date == date {
			dataList = append(dataList, a)
		}
	}
	if len(dataList) != 0 {
		for _, a := range dataTypes {
			ret = append(ret, dataAverage{Type: a, Value: getAverageValue(getValueListByType(dataList, a))})
		}
		c.IndentedJSON(http.StatusOK, gin.H{"AirportId": id, "Date": date, "TypeAverages": ret})
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
		if a.AirportID == id && a.Date == date && a.Type == dataType {
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

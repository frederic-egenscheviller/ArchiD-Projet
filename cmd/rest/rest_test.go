package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var address = "localhost:8080"

var dataArrTest = []data{
	{AirportId: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "temperature", Value: 20.0},
	{AirportId: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "wind", Value: 3.100000},
	{AirportId: "MRS", Date: "2024-01-16", Time: "12:00:00", Type: "pressure", Value: 1024.0},
	{AirportId: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "temperature", Value: 24.3},
	{AirportId: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "wind", Value: 2.0},
	{AirportId: "MRS", Date: "2024-01-16", Time: "14:15:00", Type: "pressure", Value: 1020.0},
}

func TestGetAllAirport(t *testing.T) {
	req, err := http.NewRequest("GET", "/airports", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedBody := `[
    {
        "id": "MRS"
    }
]`
	assert.Equal(t, expectedBody, rr.Body.String())
}

func TestGetAllAirportData(t *testing.T) {
	req, err := http.NewRequest("GET", "/airports/data/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAirportDataById(t *testing.T) {
	req, err := http.NewRequest("GET", "/airport/MRS/data/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetSensorsByAirportId(t *testing.T) {
	req, err := http.NewRequest("GET", "/airport/MRS/sensors", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAirportDataByDateRangesAndType(t *testing.T) {
	req, err := http.NewRequest("GET", "/airport/MRS/data/range/2024-01-16/2024-01-19/temperature", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAirportDataAverageByDate(t *testing.T) {
	req, err := http.NewRequest("GET", "/airport/MRS/average/2024-01-16", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAirportDataAverageByDateAndType(t *testing.T) {
	req, err := http.NewRequest("GET", "/airport/MRS/average/2024-01-16/temperature", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	setDataArray(dataArrTest)
	router.GET("/airports", getAllAirport)
	router.GET("/airports/data/", getAllAirportData)
	router.GET("/airport/:id/data/", getAirportDataByIATA)
	router.GET("/airport/:id/sensors", getSensorsByAirportIATA)
	router.GET("/airport/:id/data/range/:start/:end/:type", getAirportDataByDateRangesAndType)
	router.GET("/airport/:id/average/:date", getAirportDataAverageByDate)
	router.GET("/airport/:id/average/:date/:type", getAirportDataAverageByDateAndType)

	return router
}

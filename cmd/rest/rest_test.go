package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	req, err := http.NewRequest("GET", "/airport/MRS/average/2024-01-19", nil)
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
	req, err := http.NewRequest("GET", "/airport/MRS/average/2024-01-19/temperature", nil)
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

	router.GET("/airports", getAllAirports)
	router.GET("/airports/data/", getAllAirportsData)
	router.GET("/airport/:iata/data/", getAirportDataByIATA)
	router.GET("/airport/:iata/sensors", getSensorsByAirportIATA)
	router.GET("/airport/:iata/data/range/:start/:end/:measurement", getAirportDataByDateRangesAndType)
	router.GET("/airport/:iata/average/:date", getAirportDataAverageByDate)
	router.GET("/airport/:iata/average/:date/:measurement", getAirportDataAverageByDateAndType)

	return router
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testBikePointName string = "Bank of England Museum, Bank"

const testBikePointByName = `[
  {
    "id": "BikePoints_340",
    "url": "/Place/BikePoints_340",
    "commonName": "Bank of England Museum, Bank",
    "placeType": "BikePoint",
    "additionalProperties": [],
    "children": [],
    "childrenUrls": [],
    "lat": 51.514441,
    "lon": -0.087587
  }
]`

const testBikePointByNameExpected = `{
  "id": "BikePoints_340",
  "url": "/Place/BikePoints_340",
  "distance": 0,
  "commonName": "Bank of England Museum, Bank",
  "placeType": "BikePoint",
  "additionalProperties": [],
  "lat": 51.514441,
  "lon": -0.087587
}`

func TestGetBikePointByName(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	// Create new bike point object
	var bp = new(bikePoint)
	serviceEndpoint = srv.URL

	// Get the specifc bike point by name
	err := bp.getBikePointByName(testBikePointName)
	if err != nil {
		t.Errorf("Error getting bikepoint by name: %v", err)
	}

	bpjson, err := json.MarshalIndent(bp, "", "  ")
	if err != nil {
		t.Errorf("Error getting marshaling bike point: %v", err)
	}

	if testBikePointByNameExpected != string(bpjson) {
		t.Errorf("expected %s got: %s", testBikePointByNameExpected, string(bpjson))
	}
}

const testGetNumberAvailableBikes = `{
  "id": "BikePoints_340",
  "url": "/Place/BikePoints_340",
  "commonName": "Bank of England Museum, Bank",
  "placeType": "BikePoint",
  "additionalProperties": [
    {
      "category": "Description",
      "key": "NbBikes",
      "sourceSystemKey": "BikePoints",
      "value": "2",
      "modified": "2020-05-28T11:12:27.847Z"
    }
  ],
  "children": [],
  "lat": 51.514441,
  "lon": -0.087587
}
`

func TestGetNumberAvailableBikes(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	// Create new bike point object
	var bp = new(bikePoint)
	serviceEndpoint = srv.URL
	bp.ID = "BikePoints_340"

	// Get the number of bikes available at the bike point
	numberBikes, err := bp.getNumberAvailableBikes()
	if err != nil {
		t.Errorf("Error getting number of bikes: %v", err)
	}

	if 2 != numberBikes {
		t.Error("expected numberBike = 2 got", string(numberBikes))
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", rootMock)
	handler.HandleFunc("/BikePoint/Search", bikePointSearchMock)
	handler.HandleFunc("/BikePoint/BikePoints_340", bikePointLookupMock)

	srv := httptest.NewServer(handler)

	return srv
}

func bikePointSearchMock(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("query") == testBikePointName {
		_, _ = w.Write([]byte(testBikePointByName))
	} else {
		w.WriteHeader(404)
	}
}

func bikePointLookupMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(testGetNumberAvailableBikes))
}

func rootMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestRootEndpoint(t *testing.T) {
	// Test against mock api
	srv := serverMock()
	defer srv.Close()

	serviceEndpoint = srv.URL

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Welcome to the Bank of England Museum, Bank bike point checker!
There is currently 2 bike(s) at bike point: 'BikePoints_340'
which is located at: Bank of England Museum, Bank
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestHealthLivenessEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthLivenessHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `status: healthy`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestHealthReadinessEndpointGoodService(t *testing.T) {
	// Test against mock api
	srv := serverMock()
	defer srv.Close()

	serviceEndpoint = srv.URL

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthReadinessHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `status: healthy`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

func TestHealthReadinessEndpointBadService(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthReadinessHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := `status: unhealthy`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got\n%v want\n%v",
			rr.Body.String(), expected)
	}
}

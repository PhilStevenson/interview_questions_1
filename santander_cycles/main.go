package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const bikePointName string = "Bank of England Museum, Bank"

type bikePoint struct {
	ID                   string              `json:"id"`
	URL                  string              `json:"url"`
	Distance             float32             `json:"distance"`
	CommonName           string              `json:"commonName"`
	PlaceType            string              `json:"placeType"`
	AdditionalProperties []bikePointAddProps `json:"additionalProperties"`
	Lat                  float32             `json:"lat"`
	Lon                  float32             `json:"lon"`
}

type bikePointAddProps struct {
	Category        string `json:"category"`
	Key             string `json:"key"`
	SourceSystemKey string `json:"sourceSystemKey"`
	Value           string `json:"value"`
	Modified        string `json:"modified"`
}

// Gets the basic information of a single bike point by name.
func (bp *bikePoint) getBikePointByName(bikePointName string) error {
	var bikePointQuery = url.QueryEscape(bikePointName)
	var bikePointGetURL = fmt.Sprintf("https://api.tfl.gov.uk/BikePoint/Search?query=%s", bikePointQuery)

	response, err := http.Get(bikePointGetURL)
	if err != nil {
		return fmt.Errorf("The HTTP request failed with error %s", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	var bikePoints []bikePoint

	err = json.Unmarshal(data, &bikePoints)
	if err != nil {
		return fmt.Errorf("Error unmarshaling data: %s", err)
	}

	if len(bikePoints) > 1 {
		return fmt.Errorf("More than one bike point was returned using that search query")
	} else if len(bikePoints) == 0 {
		return fmt.Errorf("No bike points found")
	}

	*bp = bikePoints[0]

	return nil
}

// Gets the number of bike available at the relivant bike point. Requests additional information for the bike point.
func (bp *bikePoint) getNumberAvailableBikes() (int, error) {

	if bp.ID == "" {
		return 0, fmt.Errorf("bike point ID not set")
	}

	var bikePointID = url.QueryEscape(bp.ID)
	var bikePointGetURL = fmt.Sprintf("https://api.tfl.gov.uk/BikePoint/%s", bikePointID)

	response, err := http.Get(bikePointGetURL)
	if err != nil {
		return 0, fmt.Errorf("The HTTP request to get bike point by ID failed with error %s", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(data, bp)
	if err != nil {
		return 0, fmt.Errorf("Error unmarshaling data: %s", err)
	}

	for _, prop := range bp.AdditionalProperties {
		if prop.Key == "NbBikes" {
			i, err := strconv.Atoi(prop.Value)
			if err != nil {
				return 0, fmt.Errorf("unable to convert number of bikes to int: %v", err)
			}
			return i, nil
		}
	}

	return 0, fmt.Errorf("Unable to find the number of bikes available at: %s", bp.CommonName)
}

func main() {
	// Create new bike point object
	var bp = new(bikePoint)

	// Get the specifc bike point by name
	err := bp.getBikePointByName(bikePointName)
	if err != nil {
		fmt.Printf("Error getting bikepoint by name: %v", err)
	}

	// Get the number of bikes available at the bike point
	numberBikes, err := bp.getNumberAvailableBikes()
	if err != nil {
		fmt.Printf("Error getting number of bikes: %v", err)
	}

	// Display the information to the user
	fmt.Printf("There is currently %v bike(s) at bike point: '%s'\nwhich is located at: %s\n", numberBikes, bp.ID, bp.CommonName)
}

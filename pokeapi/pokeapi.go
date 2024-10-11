package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func GetLocationAreas(offset int) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=20&offset=%d", offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var locResponse LocationAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&locResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	var locationNames []string
	for _, location := range locResponse.Results {
		locationNames = append(locationNames, location.Name)
	}
	return locationNames, nil
}

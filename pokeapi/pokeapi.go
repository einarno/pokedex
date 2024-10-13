package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	cache "github.com/einarno/pokedexcli/pokecache"
)

type LocationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func GetLocationAreas(offset int, c *cache.Cache) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=20&offset=%d", offset)
	// Check if the key is in cache
	if cachedValue, found := c.Get(url); found {
		// Cached value is found, decode it from the cache
		var locationNames []string
		if err := json.Unmarshal(cachedValue, &locationNames); err == nil {
			// If it the decoding works we want to return, if not it's better to keep goind than returning an error
			return locationNames, nil
		}

	}

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
	encodedData, err := json.Marshal(locationNames)
	if err == nil {
		// We don't want to fail if the cache setting fails, would typically add some logging
		c.Add(url, encodedData)
	}
	return locationNames, nil
}

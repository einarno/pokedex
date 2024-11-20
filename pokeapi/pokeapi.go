package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	cache "github.com/einarno/pokedexcli/pokecache"
	pokedata "github.com/einarno/pokedexcli/pokedata"
)

type LocationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

type LocationArea struct {
	Encounters []PokemonEncounter `json:"pokemon_encounters"`
}

func GetPokemon(pokemonName string) (pokedata.Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	resp, err := http.Get(url)
	var pokemon pokedata.Pokemon
	if err != nil {
		return pokemon, fmt.Errorf("failed to fetch Pokémon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pokemon, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		return pokemon, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return pokemon, nil
}

func ExploreArea(c *cache.Cache, areaId string) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaId)
	// Check if the key is in cache
	if cachedValue, found := c.Get(url); found {
		// Cached value is found, decode it from the cache
		var pokemons []string
		if err := json.Unmarshal(cachedValue, &pokemons); err == nil {
			// If it the decoding works we want to return, if not it's better to keep goind than returning an error
			return pokemons, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location area: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	var locationArea LocationArea
	err = json.NewDecoder(resp.Body).Decode(&locationArea)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var pokemons []string
	for _, encounter := range locationArea.Encounters {
		pokemons = append(pokemons, encounter.Pokemon.Name)
	}

	return pokemons, nil
}

func GetLocationAreas(c *cache.Cache, offset int) ([]string, error) {
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

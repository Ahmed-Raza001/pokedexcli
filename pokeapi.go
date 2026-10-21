package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

const locationAreaURL = "https://pokeapi.co/api/v2/location-area"

type locationAreaPokemonResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type locationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var pokeCache = pokecache.NewCache(5 * time.Second)

func getLocationAreas(url string) (locationAreasResponse, error) {
	var result locationAreasResponse

	data, ok := pokeCache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return result, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return result, fmt.Errorf("API returned status: %s", resp.Status)
		}

		var responseBytes []byte
		buffer := make([]byte, 4096)

		for {
			n, readErr := resp.Body.Read(buffer)
			if n > 0 {
				responseBytes = append(responseBytes, buffer[:n]...)
			}
			if readErr != nil {
				if readErr.Error() != "EOF" {
					return result, readErr
				}
				break
			}
		}

		data = responseBytes
		pokeCache.Add(url, data)
	}

	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil

}

func getLocationAreaPokemon(url string) (locationAreaPokemonResponse, error) {
	var result locationAreaPokemonResponse

	data, ok := pokeCache.Get(url)

	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return result, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return result, fmt.Errorf("API returned status: %s", resp.Status)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}

		pokeCache.Add(url, data)
	}

	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func getPokemon(name string) (Pokemon, error) {
	var pokemon Pokemon
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	data, ok := pokeCache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return pokemon, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return pokemon, fmt.Errorf("pokemon not found: %s", name)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return pokemon, err
		}

		pokeCache.Add(url, data)
	}

	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return pokemon, err
	}

	return pokemon, nil

}

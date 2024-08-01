package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nrm-mrn/pokedexcli/internal/pokecache"
)

type LocationsResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const baseURL = "https://pokeapi.co/api/v2"

const staleTime = time.Second * 5

var cache = pokecache.NewCache(staleTime)

func getLocationsStruct(responseData []byte) (LocationsResponse, error) {
	locations := LocationsResponse{}
	err := json.Unmarshal(responseData, &locations)
	if err != nil {
		return locations, err
	}
	return locations, nil
}

func GetLocations(pageUrl *string) (LocationsResponse, error) {
	endpointURL := "/location-area"
	fullUrl := baseURL + endpointURL
	if pageUrl != nil {
		fullUrl = *pageUrl
	}

	locations := LocationsResponse{}
	var body []byte
	cached, ok := cache.Get(fullUrl)
	if ok {
		body = cached
	} else {
		res, err := http.Get(fullUrl)
		defer res.Body.Close()
		if err != nil {
			return locations, err
		}
		if res.StatusCode > 399 {
			errText := fmt.Sprintf(
				"Response failed with status code %d and\nbody: %v\n",
				res.StatusCode,
				res.Body,
			)
			return locations, errors.New(errText)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return locations, err
		}
	}

	result, err := getLocationsStruct(body)
	if err != nil {
		return locations, err
	}
	cache.Add(fullUrl, body)

	return result, nil
}

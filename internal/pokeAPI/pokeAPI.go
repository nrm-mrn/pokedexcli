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

const baseURL = "https://pokeapi.co/api/v2"

const staleTime = time.Second * 10

var cache = pokecache.NewCache(staleTime)

func getLocationsStruct(responseData []byte) (LocationsResponse, error) {
	locations := LocationsResponse{}
	err := json.Unmarshal(responseData, &locations)
	if err != nil {
		return locations, err
	}
	return locations, nil
}

func getReq(URL string) ([]byte, error) {
	res, err := http.Get(URL)
	defer res.Body.Close()
	if err != nil {
		return []byte(""), err
	}
	if res.StatusCode > 399 {
		errText := fmt.Sprintf(
			"Response failed with status code %d and \nbody: %v",
			res.StatusCode,
			res.Body,
		)
		return []byte(""), errors.New(errText)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte(""), err
	}
	return body, nil
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
		var err error
		body, err = getReq(fullUrl)
		if err != nil {
			return locations, err
		}
	}

	result, err := getLocationsStruct(body)
	if err != nil {
		return locations, err
	}
	if !ok {
		cache.Add(fullUrl, body)
	}

	return result, nil
}

func getExploredStruct(respBody []byte) (ExploreResponse, error) {
	var exploredLoca ExploreResponse
	err := json.Unmarshal(respBody, &exploredLoca)
	if err != nil {
		return exploredLoca, err
	}
	return exploredLoca, nil
}
func ExploreLocation(pageUrl *string) (ExploreResponse, error) {
	endpointURL := baseURL + "/location-area/" + *pageUrl
	exploredLoca := ExploreResponse{}
	if pageUrl == nil {
		return exploredLoca, errors.New("Nil location name, use map for list of all names")
	}

	cached, ok := cache.Get(endpointURL)
	var body []byte
	var err error
	if ok {
		body = cached
	} else {
		body, err = getReq(endpointURL)
		if err != nil {
			return exploredLoca, err
		}
	}
	exploredLoca, err = getExploredStruct(body)
	if err != nil {
		return exploredLoca, err
	}
	if !ok {
		cache.Add(endpointURL, body)
	}
	return exploredLoca, nil
}

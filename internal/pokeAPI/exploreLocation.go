package pokeapi

import (
	"encoding/json"
	"errors"
)

func getExploredStruct(respBody []byte) (ExploreResponse, error) {
	var exploredLoca ExploreResponse
	err := json.Unmarshal(respBody, &exploredLoca)
	if err != nil {
		return exploredLoca, err
	}
	return exploredLoca, nil
}

func ExploreLocation(c *Client, pageUrl *string) (ExploreResponse, error) {
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
		body, err = getReq(c, endpointURL)
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

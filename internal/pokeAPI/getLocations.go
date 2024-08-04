package pokeapi

import "encoding/json"

func getLocationsStruct(responseData []byte) (LocationsResponse, error) {
	locations := LocationsResponse{}
	err := json.Unmarshal(responseData, &locations)
	if err != nil {
		return locations, err
	}
	return locations, nil
}

func GetLocations(c *Client, pageUrl *string) (LocationsResponse, error) {
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
		body, err = getReq(c, fullUrl)
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

package pokeapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nrm-mrn/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

const staleTime = time.Second * 10

type Client struct {
	httpClient http.Client
}

var cache = pokecache.NewCache(staleTime)

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
	}
}
func getReq(c *Client, URL string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return []byte(""), err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer res.Body.Close()

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

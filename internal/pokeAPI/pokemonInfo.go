package pokeapi

import "encoding/json"

func getPokemonStruct(body []byte, pokemonInfo *Pokemon) error {
	err := json.Unmarshal(body, pokemonInfo)
	return err

}

func GetPokemon(c *Client, pokeName string) (Pokemon, error) {
	fullURL := baseURL + "/pokemon/" + pokeName
	var pokemonInfo Pokemon
	var body []byte
	var err error

	cached, ok := cache.Get(fullURL)
	if ok {
		body = cached
	} else {
		body, err = getReq(c, fullURL)
		if err != nil {
			return pokemonInfo, err
		}
	}
	err = getPokemonStruct(body, &pokemonInfo)
	if err != nil {
		return pokemonInfo, err
	}
	if !ok {
		cache.Add(fullURL, body)
	}
	return pokemonInfo, nil
}

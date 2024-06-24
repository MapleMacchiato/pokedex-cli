package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MapleMacchiato/pokedex-cli/internal/pokecache"
)

type Client struct {
	locationsCache pokecache.Cache
	areaCache      pokecache.Cache
	PrevURL        *string
	NextURL        *string
	httpClient     http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		locationsCache: pokecache.NewCache(cacheInterval),
		areaCache:      pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetLocations(pageURL *string) error {
	var url *string
	if pageURL == nil {
		temp := "https://pokeapi.co/api/v2/location-area/"
		url = &temp
	} else {
		url = pageURL
	}

	locations := LocationAreas{}
	body, ok := c.locationsCache.Get(*url)
	if ok {
		err := json.Unmarshal(body, &locations)
		if err != nil {
			fmt.Println("test")
			return err
		}
	} else {

		req, err := http.NewRequest("GET", *url, nil)
		if err != nil {
			return err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		c.locationsCache.Add(*url, body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &locations)
		if err != nil {
			return err
		}
	}

	c.NextURL = locations.Next
	c.PrevURL = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func (c *Client) GetPokemon(area string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + area
	pokemons := Location{}
	body, ok := c.areaCache.Get(url)
	if ok {
		err := json.Unmarshal(body, &pokemons)
		if err != nil {
			return err
		}
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		c.areaCache.Add(url, body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &pokemons)
		if err != nil {
			return err
		}
	}
	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

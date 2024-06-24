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

type Locations struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Pokemons struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	}
}

func (l *Locations) printLocations() {
	var names string
	for _, location := range *&l.Results {
		fmt.Println(location.Name)
		names += location.Name + "\n"
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

	locations := Locations{}
	body, ok := c.locationsCache.Get(*url)
	if ok {
		err := json.Unmarshal(body, &locations)
		if err != nil {
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
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &locations)
		if err != nil {
			return err
		}
	}
	c.locationsCache.Add(*url, body)
	c.NextURL = locations.Next
	c.PrevURL = locations.Previous
	locations.printLocations()
	return nil
}

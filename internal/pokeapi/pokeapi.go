package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/MapleMacchiato/pokedex-cli/internal/pokecache"
)

type Client struct {
	Pokedex    map[string]Pokemon
	cache      pokecache.Cache
	PrevURL    *string
	NextURL    *string
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		Pokedex: make(map[string]Pokemon),
		cache:   pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) PrintPokedex() error {
	fmt.Println("Your Pokedex: ")
	for k := range c.Pokedex {
		fmt.Printf(" - %s\n", k)
	}
	return nil
}

func (c *Client) Inspect(pokemon string) error {
	if pk, ok := c.Pokedex[pokemon]; ok {
		fmt.Println("Name:", pk.Name)
		fmt.Println("Height:", pk.Height)
		fmt.Println("Weight:", pk.Weight)
		fmt.Println("Stats:")
		for _, stat := range pk.Stats {
			fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeInfo := range pk.Types {
			fmt.Println("  -", typeInfo.Type.Name)
		}
		return nil
	}
	return errors.New("you do not have that pokemon")

}

func (c *Client) GetLocations(pageURL *string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pageURL != nil {
		url = string(*pageURL)
	}

	locations := LocationAreas{}
	err := c.makeRequest(url, &locations)
	if err != nil {
		return err
	}

	c.NextURL = locations.Next
	c.PrevURL = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func (c *Client) ExploreArea(area string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + area
	pokemons := Location{}
	err := c.makeRequest(url, &pokemons)
	if err != nil {
		return err
	}
	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func (c *Client) CatchPokemon(pokemon string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	pk := Pokemon{}
	err := c.makeRequest(url, &pk)
	if err != nil {
		return err
	}
	catchValue := rand.Intn(10)
	fmt.Printf("Throwing pokeball at %s...\n", pokemon)
	if catchValue > 4 {
		fmt.Println("Success!")
		c.Pokedex[pokemon] = pk
		return nil
	}
	fmt.Printf("Failed to catch %s\n", pokemon)
	return nil
}

func (c *Client) makeRequest(url string, s any) error {
	body, ok := c.cache.Get(url)
	if ok {
		err := json.Unmarshal(body, &s)
		if err != nil {
			return err
		}
		return nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	c.cache.Add(url, body)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return err
	}
	return nil
}

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

func (c *Client) Inspect(pokemonName string) error {
	if pokemon, ok := c.Pokedex[pokemonName]; ok {
		fmt.Println("Name:", pokemon.Name)
		fmt.Println("Height:", pokemon.Height)
		fmt.Println("Weight:", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeInfo := range pokemon.Types {
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

func (c *Client) CatchPokemon(pokemonName string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	pokemon := Pokemon{}
	err := c.makeRequest(url, &pokemon)
	if err != nil {
		return err
	}
	catchValue := rand.Intn(10)
	fmt.Printf("Throwing pokeball at %s...\n", pokemonName)
	if catchValue > 4 {
		fmt.Println("Success!")
		c.Pokedex[pokemonName] = pokemon
		return nil
	}
	fmt.Printf("Failed to catch %s\n", pokemonName)
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

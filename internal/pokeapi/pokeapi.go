package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Areas struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

var areas = &Areas{}

func (a *Areas) getMaps() {
	var url string
	if a.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *a.Next
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &a)
	if err != nil {
		fmt.Println(err)
	}

	for _, location := range *&a.Results {
		fmt.Println(location.Name)
	}
}

func (a *Areas) getMapsB() {
	var url string
	if a.Previous == nil {
		fmt.Println("No previous areas")
		return
	} else {
		url = *a.Previous
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &a)
	if err != nil {
		fmt.Println(err)
	}

	for _, location := range *&a.Results {
		fmt.Println(location.Name)
	}

}

func GetMapsB() {
	areas.getMapsB()
}

func GetMaps() {
	areas.getMaps()
}

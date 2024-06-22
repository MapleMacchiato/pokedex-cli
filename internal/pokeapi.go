package internal

import (
	"fmt"
	"io"
	"net/http"
)

func GetMaps() {
	url := "https://pokeapi.co/api/v2/location-area/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error")
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error")
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(res)
	fmt.Println(string(body))
}

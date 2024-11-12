package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type locationResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(url string) (*string, string, []string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, "", nil, err
	}

	locations := locationResponse{}
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&locations); err != nil {
		return nil, "", nil, err
	}

	locationNames := make([]string, len(locations.Results))
	for i := 0; i < len(locationNames); i++ {
		locationNames[i] = locations.Results[i].Name
	}
	fmt.Println(locationNames)
	if locations.Previous != nil {
		return locations.Previous, locations.Next, locationNames, nil
	} else {
		return nil, locations.Next, locationNames, nil
	}

}

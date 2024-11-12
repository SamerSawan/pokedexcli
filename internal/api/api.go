package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/samersawan/pokedexcli/internal/pokecache"
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

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (client *Client) GetLocations(url string, c *pokecache.Cache) (*string, string, []string, error) {
	locations := locationResponse{}
	if dat, exists := c.Get(url); exists {
		err := json.Unmarshal(dat, &locations)
		if err != nil {
			return nil, "", nil, err
		}
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, "", nil, err
		}

		res, err := client.httpClient.Do(req)
		if err != nil {
			return nil, "", nil, err
		}
		defer res.Body.Close()

		dat, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, "", nil, err
		}

		c.Add(url, dat)

		err = json.Unmarshal(dat, &locations)
		if err != nil {
			return nil, "", nil, err
		}
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

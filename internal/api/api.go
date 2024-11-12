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

type pokemonResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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

func (client *Client) ExploreLocation(url string, c *pokecache.Cache) ([]string, error) {
	pokemon := pokemonResponse{}
	if dat, exists := c.Get(url); exists {
		err := json.Unmarshal(dat, &pokemon)
		if err != nil {
			return nil, err
		}
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		res, err := client.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		dat, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		c.Add(url, dat)

		err = json.Unmarshal(dat, &pokemon)
		if err != nil {
			return nil, err
		}
	}

	pokemonNames := make([]string, len(pokemon.PokemonEncounters))
	for i := 0; i < len(pokemonNames); i++ {
		pokemonNames[i] = pokemon.PokemonEncounters[i].Pokemon.Name
	}
	return pokemonNames, nil
}

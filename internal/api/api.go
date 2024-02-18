package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/unclesamrocks/pokedexcli/internal/http"
	"github.com/unclesamrocks/pokedexcli/internal/pokecache"
)

type apiPokemon struct {
	next  *string
	prev  *string
	cache *pokecache.Cache
}

func New() *apiPokemon {
	ctx := apiPokemon{
		next:  nil,
		prev:  nil,
		cache: pokecache.NewCache(time.Second * 60),
	}

	return &ctx
}

func (apiPokemon *apiPokemon) FetchNext(args ...string) error {
	url := BASE_URL + LOCATION_AREAS_URL

	if apiPokemon.next != nil {
		url = *apiPokemon.next
	}

	data, errCache := apiPokemon.cache.Get(url)

	if !errCache {
		fmt.Println("Fetching...")
		fetchData, errGet := http.Get(url)

		if errGet != nil {
			return errGet
		}

		data = fetchData
		apiPokemon.cache.Add(url, data)
	}

	locationAreas := LocationAreas{}

	if errJson := json.Unmarshal(data, &locationAreas); errJson != nil {
		return errJson
	}

	apiPokemon.next = locationAreas.Next
	apiPokemon.prev = locationAreas.Previous

	printLocationAreas(&locationAreas)

	return nil
}

func (apiPokemon *apiPokemon) FetchPrev(args ...string) error {
	if apiPokemon.prev == nil {
		error := errors.New("no prev page available")
		return error
	}

	url := *apiPokemon.prev

	data, errCache := apiPokemon.cache.Get(url)

	if !errCache {
		fmt.Println("Fetching...")
		fetchData, errGet := http.Get(url)

		if errGet != nil {
			return errGet
		}

		data = fetchData
		apiPokemon.cache.Add(url, data)
	}

	locationAreas := LocationAreas{}

	if errJson := json.Unmarshal(data, &locationAreas); errJson != nil {
		return errJson
	}

	apiPokemon.next = locationAreas.Next
	apiPokemon.prev = locationAreas.Previous

	printLocationAreas(&locationAreas)

	return nil
}

func printLocationAreas(locationAreas *LocationAreas) {
	for _, location := range (*locationAreas).Results {
		fmt.Printf("%s [%s]\n", location.Name, location.URL)
	}
}

func (c *apiPokemon) Explore(args ...string) error {
	locationId := args[1]
	baseUrl, err := url.Parse(BASE_URL + LOCATION_AREAS_URL)

	if err != nil {
		return err
	}

	locationUrl := baseUrl.JoinPath(locationId).String()

	rawData, errCacheData := c.cache.Get(locationUrl)

	if !errCacheData {
		fetchData, errFetchData := http.Get(locationUrl)

		if errFetchData != nil {
			return errFetchData
		}

		rawData = fetchData
		c.cache.Add(locationUrl, fetchData)
	}

	data := Area{}

	if errJson := json.Unmarshal(rawData, &data); errJson != nil {
		return errJson
	}

	// prettyJSON, err := json.MarshalIndent(&data, "", "  ")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(prettyJSON))

	for _, v := range data.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}

	return nil
}

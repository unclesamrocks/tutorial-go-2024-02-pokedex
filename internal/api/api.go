package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/unclesamrocks/pokedexcli/internal/http"
	"github.com/unclesamrocks/pokedexcli/internal/pokecache"
)

type apiPokemon struct {
	next    *string
	prev    *string
	cache   *pokecache.Cache
	pokedex map[string]Pokemon
}

func New() *apiPokemon {
	ctx := apiPokemon{
		next:    nil,
		prev:    nil,
		cache:   pokecache.NewCache(time.Second * 60),
		pokedex: map[string]Pokemon{},
	}

	return &ctx
}

func (a *apiPokemon) FetchNext(args ...string) error {
	url := BASE_URL + LOCATION_AREAS_URL

	if a.next != nil {
		url = *a.next
	}

	fmt.Println(url)
	rawData, err := a.fetch(url)

	if err != nil {
		return err
	}

	locationAreas := LocationAreas{}

	if err := json.Unmarshal(rawData, &locationAreas); err != nil {
		return err
	}

	fmt.Printf("Next: %s\n", *locationAreas.Next)

	a.next = locationAreas.Next
	a.prev = locationAreas.Previous

	printLocationAreas(&locationAreas)

	return nil
}

func (a *apiPokemon) FetchPrev(args ...string) error {
	if a.prev == nil {
		err := errors.New("no prev page available")
		return err
	}

	url := *a.prev

	rawData, err := a.fetch(url)

	if err != nil {
		return err
	}

	locationAreas := LocationAreas{}

	if err := json.Unmarshal(rawData, &locationAreas); err != nil {
		return err
	}

	a.next = locationAreas.Next
	a.prev = locationAreas.Previous

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

	rawData, err := c.fetch(locationUrl)

	if err != nil {
		return err
	}

	data := Area{}

	if errJson := json.Unmarshal(rawData, &data); errJson != nil {
		return errJson
	}

	for _, v := range data.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}

	return nil
}

func (c *apiPokemon) Catch(args ...string) error {
	pokemonId := args[1]
	baseUrl, err := url.Parse(BASE_URL + POKEMON)

	if err != nil {
		return err
	}

	pokemonUrl := baseUrl.JoinPath(pokemonId).String()

	rawData, err := c.fetch(pokemonUrl)

	if err != nil {
		return err
	}

	data := Pokemon{}

	if errJson := json.Unmarshal(rawData, &data); errJson != nil {
		return errJson
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonId)
	if rand.Intn(2) == 1 {
		fmt.Printf("%s was caught!\n", pokemonId)
		c.pokedex[pokemonId] = data
	} else {
		fmt.Printf("%s escaped!\n", pokemonId)
	}

	return nil
}

func (a *apiPokemon) fetch(url string) ([]byte, error) {
	rawData, errCacheData := a.cache.Get(url)

	if !errCacheData {
		fmt.Println("Fetching...")
		fetchData, errFetchData := http.Get(url)

		if errFetchData != nil {
			return []byte{}, errFetchData
		}

		rawData = fetchData
		a.cache.Add(url, fetchData)
	}

	return rawData, nil
}

func (a *apiPokemon) prettyPrint(data any) error {
	prettyJSON, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON))

	return nil
}

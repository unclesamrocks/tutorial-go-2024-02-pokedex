package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/unclesamrocks/pokedexcli/internal/http"
)

type apiPokemon struct {
	next *string
	prev *string
}

func New() *apiPokemon {
	ctx := apiPokemon{
		next: nil,
		prev: nil,
	}

	return &ctx
}

func (apiPokemon *apiPokemon) FetchNext() error {
	url := BASE_URL + LOCATION_AREAS_URL

	if apiPokemon.next != nil {
		url = *apiPokemon.next
	}

	data, errGet := http.Get(url)

	if errGet != nil {
		return errGet
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

func (apiPokemon *apiPokemon) FetchPrev() error {
	if apiPokemon.prev == nil {
		error := errors.New("no prev page available")
		return error
	}

	body, err := http.Get(*apiPokemon.prev)

	if err != nil {
		return err
	}

	locationAreas := LocationAreas{}

	if errJson := json.Unmarshal(body, &locationAreas); errJson != nil {
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

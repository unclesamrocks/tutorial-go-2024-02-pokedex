package api

type LocationAreas struct {
	Count    int                    `json:"count"`
	Next     *string                `json:"next"`
	Previous *string                `json:"previous"`
	Results  []LocationAreasResults `json:"results"`
}

type LocationAreasResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

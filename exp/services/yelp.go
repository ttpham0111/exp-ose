package services

import (
	"net/http"
	"net/url"
)

const businessAPI = "/businesses"

type YelpService struct {
	apiURL string
	apiKey string
}

func NewYelpService(apiKey string) *YelpService {
	return &YelpService{
		apiURL: "https://api.yelp.com/v3",
		apiKey: apiKey,
	}
}

type BusinessQuery struct {
	term      string
	location  string
	latitude  string
	longitude string
}

func (q *BusinessQuery) rawQuery() string {
	payload = url.Values{}

	payload.add("term", q.term)

	if q.location != "" {
		payload.add("location", q.location)
	}

	if q.latitude != "" {
		payload.add("latitude", q.latitude)
	}

	if q.longitude != "" {
		payload.add("longitude", q.longitude)
	}

	if q.limit > 0 {
		payload.add("limit", q.limit)
	}

	return payload.Encode()
}

type Location struct {
	city     string
	country  string
	address2 string
	address3 string
	state    string
	address1 string
	zipCode  string `json: zip_code`
}

type Business struct {
	name        string
	url         string
	imageURL    string
	rating      int
	reviewCount int `json: review_count`
	location    Location
}

func (service *YelpService) findBusinesses(query *BusinessQuery) (businesses []Business) {
	res, err := http.Get(service.apiURL + businessAPI + "/search?" + payload.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = raiseForStatus(res); err != nil {
		return nil, err
	}

	var body map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nul, err
	}

	businesses = body["businesses"].([]Business)
}

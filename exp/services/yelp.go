package services

import (
	"encoding/json"
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

type YelpLocation struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	State    string `json:"state"`
	ZipCode  string `json: zip_code`
}

type YelpBusiness struct {
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	ImageURL    string       `json:"image_url"`
	Rating      float32      `json:"rating"`
	ReviewCount int          `json:"review_count"`
	Location    YelpLocation `json:"location"`
}

func (service *YelpService) FindBusinesses(query Query) ([]YelpBusiness, error) {
	url := service.apiURL + businessAPI + "/search?" + query.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+service.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = raiseForStatus(res); err != nil {
		return nil, err
	}

	var body struct {
		Businesses []YelpBusiness `json:"businesses"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Businesses, nil
}

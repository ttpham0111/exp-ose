package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const businessesAPI = "/businesses"

type YelpService interface {
	FindBusinesses(ServiceQuery) ([]*YelpBusiness, error)
}

type yelpServiceContext struct {
	apiURL string
	apiKey string
}

func NewYelpService(apiKey string) *yelpServiceContext {
	return &yelpServiceContext{
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
	ZipCode  string `json:"zip_code"`
}

type YelpBusiness struct {
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	ImageURL    string       `json:"image_url"`
	Rating      float32      `json:"rating"`
	ReviewCount int          `json:"review_count"`
	Location    YelpLocation `json:"location"`
}

func encodeYelpQuery(q ServiceQuery) string {
	v := url.Values{}
	v.Add("term", q.Term)

	if q.Location != "" {
		v.Add("location", q.Location)
	}

	if q.Latitude != "" && q.Longitude != "" {
		v.Add("latitude", q.Latitude)
		v.Add("longitude", q.Longitude)
	}

	if q.Limit > 0 {
		v.Add("limit", strconv.Itoa(q.Limit))
	}

	return v.Encode()
}

func (ctx *yelpServiceContext) FindBusinesses(query ServiceQuery) ([]*YelpBusiness, error) {
	url := ctx.apiURL + businessesAPI + "/search?" + encodeYelpQuery(query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+ctx.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = raiseForStatus(res); err != nil {
		return nil, err
	}

	var body struct {
		Businesses []*YelpBusiness `json:"businesses"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Businesses, nil
}

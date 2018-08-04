package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const businessAPI = "/businesses"

type YelpService interface {
	FindBusinesses(Query) ([]YelpBusiness, error)
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

type YelpQuery struct {
	Term      string   `form:"term"`
	Location  string   `form:"location"`
	Latitude  *float64 `form:"latitude" binding:"omitempty,latitude"`
	Longitude *float64 `form:"longitude" binding:"omitempty,longitude"`
	Limit     int      `form:"limit" binding:"omitempty,gt=0"`
}

func (q YelpQuery) Encode() string {
	v := url.Values{}
	v.Add("term", q.Term)

	if q.Location != "" {
		v.Add("location", q.Location)
	}

	if q.Latitude != nil && q.Longitude != nil {
		v.Add("latitude", strconv.FormatFloat(*q.Latitude, 'f', -1, 64))
		v.Add("longitude", strconv.FormatFloat(*q.Latitude, 'f', -1, 64))
	}

	if q.Limit > 0 {
		v.Add("limit", strconv.Itoa(q.Limit))
	}

	return v.Encode()
}

func (ctx *yelpServiceContext) FindBusinesses(query Query) ([]YelpBusiness, error) {
	url := ctx.apiURL + businessAPI + "/search?" + query.Encode()
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
		Businesses []YelpBusiness `json:"businesses"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Businesses, nil
}

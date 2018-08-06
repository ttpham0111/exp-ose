package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const eventsAPI = "/events"

type EventbriteService interface {
	FindEvents(ServiceQuery) ([]*EventbriteEvent, error)
}

type eventbriteServiceContext struct {
	apiURL   string
	apiToken string
}

func NewEventbriteService(apiToken string) *eventbriteServiceContext {
	return &eventbriteServiceContext{
		apiURL:   "https://www.eventbriteapi.com/v3/",
		apiToken: apiToken,
	}
}

type MultipartText struct {
	Text string `json:"text"`
	Html string `json:"html"`
}

type DatetimeTz struct {
	Timezone string `json:"timezone"`
	Utc      string `json:"utc"`
	Local    string `json:"local"`
}

type Image struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type EventbriteEvent struct {
	Name        MultipartText `json:"name"`
	Description MultipartText `json:"description"`
	Url         string        `json:"url"`
	StartsAt    DatetimeTz    `json:"start"`
	EndsAt      DatetimeTz    `json:"end"`
	Status      string        `json:"status"`
	OnlineEvent bool          `json:"online_event"`
	Logo        Image         `json:"logo"`
}

func encodeEventbriteQuery(q ServiceQuery) string {
	v := url.Values{}
	v.Add("q", q.Term)
	v.Add("start_date.range_start", time.Now().Format("2006-01-02T15:04:05"))

	if q.Location != "" {
		v.Add("location.address", q.Location)
	}

	if q.Latitude != "" && q.Longitude != "" {
		v.Add("location.latitude", q.Latitude)
		v.Add("location.longitude", q.Longitude)
	}

	if q.Limit > 0 {
		v.Add("limit", strconv.Itoa(q.Limit))
	}

	return v.Encode()
}

func (ctx *eventbriteServiceContext) FindEvents(query ServiceQuery) ([]*EventbriteEvent, error) {
	url := ctx.apiURL + eventsAPI + "/search?" + encodeEventbriteQuery(query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+ctx.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = raiseForStatus(res); err != nil {
		return nil, err
	}

	var body struct {
		Events []*EventbriteEvent `json:"events"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Events, nil
}

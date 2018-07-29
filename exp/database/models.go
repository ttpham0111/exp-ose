package database

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
)

type EventSource int

const (
	DB EventSource = iota
	Yelp
	Eventbrite
	Google
	User
)

var eventSourceItoa = map[EventSource]string{
	DB:         "db",
	Yelp:       "yelp",
	Eventbrite: "eventbrite",
	Google:     "google",
	User:       "user",
}

var eventSourceAtoi = reverseMap(eventSourceItoa)

func reverseMap(m map[EventSource]string) map[string]EventSource {
	n := make(map[string]EventSource)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func (source *EventSource) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(eventSourceItoa[*source])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (source *EventSource) UnmarshalJSON(buffer []byte) error {
	var name string
	if err := json.Unmarshal(buffer, &name); err != nil {
		return err
	}

	*source = eventSourceAtoi[name]
	return nil
}

type Location struct {
	City    string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
	Address string `json:"address" bson:"address"`
	State   string `json:"state" bson:"state"`
	ZipCode string `json:"zip_code" bson:"zip_code"`
}

type Event struct {
	Id             bson.ObjectId `json:"id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Url            string        `json:"url" bson:"url"`
	ImageURL       string        `json:"image_url" bson:"image_url"`
	Location       Location      `json:"location" bson:"location"`
	StartsAt       time.Time     `json:"time" bson:"time"`
	EndsAt         time.Time     `json:"time" bson:"time"`
	Source         EventSource   `json:"source" bson:"source"`
	SourceMetadata interface{}   `json:"source_metadata" bson:"source_metadata"`
}

type Experience struct {
	Id       bson.ObjectId   `json:"id" bson:"_id"`
	Owner    string          `json:"owner" bson:"owner"`
	Name     string          `json:"name" bson:"name"`
	ImageURL string          `json:"image_url" bson:"image_url"`
	Events   []bson.ObjectId `json:"events" bson:"events"`
	Tags     []string        `json:"tags" bson:"tags"`
	IsPublic bool            `json:"is_public" bson:"is_public"`
}

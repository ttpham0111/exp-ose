package database

import (
	"bytes"
	"encoding/json"
	"net/url"
	"time"

	"github.com/globalsign/mgo/bson"
)

type UserId string
type EventSource int

const (
	User = iota
	Yelp
	Eventbrite
	Google
)

var eventSourceItoa = map[EventSource]string{
	User:       "user",
	Yelp:       "yelp",
	Eventbrite: "eventbrite",
	Google:     "google",
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

type Experience struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Owner       UserId        `json:"owner" bson:"owner"`
	IsPrivate   bool          `json:"is_private" bson:"is_private"`
	Name        string        `json:"name" bson:"name"`
	ImageURL    string        `json:"image_url" bson:"image_url"`
	Rating      float32       `json:"rating" bson:"rating"`
	Tags        []string      `json:"tags" bson:"tags"`
	NumEvents   int           `json:"num_events" bson:"num_events"`
	NumComments int           `json:"num_comments" bson:"num_comments"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}

type Collaborators struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	ExperienceId  bson.ObjectId `json:"experience_id" bson:"experience_id"`
	Collaborators []UserId      `json:"collaborators" bson:"collaborators"`
}

type ExperienceEvent struct {
	Id           bson.ObjectId `json:"id" bson:"_id"`
	ExperienceId bson.ObjectId `json:"experience_id" bson:"experience_id"`
	StartsAt     time.Time     `json:"time" bson:"time"`
	EndsAt       time.Time     `json:"time" bson:"time"`
	Event
}

type Comment struct {
	Id           bson.ObjectId `json:"id" bson:"_id"`
	ExperienceId bson.ObjectId `json:"experience_id" bson:"experience_id"`
	Owner        UserId        `json:"owner" bson:"owner"`
	Text         string        `json:"text" bson:"text"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
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
	ImageURL       url.URL       `json:"image_url" bson:"image_url"`
	Source         EventSource   `json:"source" bson:"source"`
	SourceMetadata interface{}   `json:"source_metadata" bson:"source_metadata"`
}

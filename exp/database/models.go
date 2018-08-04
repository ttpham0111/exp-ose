package database

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
)

type UserId string
type ActivitySource int

const (
	User = iota
	Yelp
	Eventbrite
	Google
)

var activitySourceItoa = map[ActivitySource]string{
	User:       "user",
	Yelp:       "yelp",
	Eventbrite: "eventbrite",
	Google:     "google",
}

var ActivitySourceAtoi = reverseMap(activitySourceItoa)

func reverseMap(m map[ActivitySource]string) map[string]ActivitySource {
	n := make(map[string]ActivitySource)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func (source *ActivitySource) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(activitySourceItoa[*source])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (source *ActivitySource) UnmarshalJSON(buffer []byte) error {
	var name string
	if err := json.Unmarshal(buffer, &name); err != nil {
		return err
	}

	val, exists := ActivitySourceAtoi[name]
	if !exists {
		return ValidationError{"source"}
	}

	*source = val
	return nil
}

type Experience struct {
	Id            bson.ObjectId        `json:"id" bson:"_id"`
	Owner         UserId               `json:"owner" bson:"owner" binding:"required"`
	Collaborators []UserId             `json:"collaborators" bson:"collaborators" binding:"omitempty,gt=0"`
	IsPrivate     bool                 `json:"is_private" bson:"is_private"`
	Name          string               `json:"name" bson:"name" binding:"required"`
	ImageURL      string               `json:"image_url" bson:"image_url"`
	Rating        float32              `json:"rating" bson:"rating"`
	NumRatings    int                  `json:"num_ratings" bson:"num_ratings"`
	Tags          []string             `json:"tags" bson:"tags"`
	Activities    []ExperienceActivity `json:"activities" bson:"activities"`
}

type Comment struct {
	Id           bson.ObjectId `json:"id" bson:"_id"`
	ExperienceId bson.ObjectId `json:"experience_id" bson:"experience_id"`
	Owner        UserId        `json:"owner" bson:"owner" binding:"required"`
	Text         string        `json:"text" bson:"text" binding:"required"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
}

type SourceMetadata map[string]interface{}

type ExperienceActivity struct {
	StartsAt       time.Time      `json:"starts_at" bson:"starts_at"`
	EndsAt         time.Time      `json:"ends_at" bson:"ends_at"`
	Name           string         `json:"name" bson:"name" binding:"required"`
	ImageURL       string         `json:"image_url" bson:"image_url"`
	Source         ActivitySource `json:"source" bson:"source" binding:"exists"`
	SourceMetadata SourceMetadata `json:"source_metadata" bson:"source_metadata"`
}

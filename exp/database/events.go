package database

import (
	"github.com/globalsign/mgo"
)

const eventType = "Event"

type EventCollectionReader interface {
	Find(EventQuery) ([]*Event, error)
}

type EventCollection struct {
	c *mgo.Collection
}

type EventQuery struct {
	Name string
	Location
	Source EventSource
}

func (c *EventCollection) Find(query EventQuery) ([]*Event, error) {
	events := make([]*Event, 0)
	err := c.c.Find(query).All(&events)
	return events, err
}

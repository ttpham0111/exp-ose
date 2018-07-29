package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type EventCollectionReader interface {
	Find(EventQuery) (Event, error)
	FindIds([]bson.ObjectId) ([]Event, error)
}

type EventCollection struct {
	c *mgo.Collection
}

type EventQuery bson.M

func (c *EventCollection) Find(query EventQuery) (events Event, err error) {
	err = c.c.Find(query).All(&events)
	return events, err
}

func (c *EventCollection) FindIds(ids []bson.ObjectId) (events []Event, err error) {
	err = c.c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&events)
	return events, err
}

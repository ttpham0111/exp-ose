package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type ExperienceCollectionReader interface {
	Find(ExperienceQuery) ([]Experience, error)
	FindId(bson.ObjectId) (Experience, error)
}

type ExperienceCollection struct {
	c *mgo.Collection
}

type ExperienceQuery bson.M

func (c *ExperienceCollection) Find(query ExperienceQuery) (experiences []Experience, err error) {
	err = c.c.Find(query).All(&experiences)
	return experiences, err
}

func (c *ExperienceCollection) FindId(id bson.ObjectId) (experience Experience, err error) {
	err = c.c.FindId(id).One(&experience)
	return experience, err
}

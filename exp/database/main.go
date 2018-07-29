package database

import (
	"github.com/globalsign/mgo"
)

type Database struct {
	session *mgo.Session
	*EventCollection
	*ExperienceCollection
}

func NewDatabase(url string, dbName string) (*Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	db := session.DB(dbName)

	return &Database{
		session:              session,
		EventCollection:      &EventCollection{db.C("events")},
		ExperienceCollection: &ExperienceCollection{db.C("experiences")},
	}, nil
}

func (db *Database) Close() {
	db.session.Close()
}

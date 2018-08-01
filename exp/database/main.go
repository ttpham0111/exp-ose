package database

import (
	"github.com/globalsign/mgo"
)

const (
	experienceCollection      = "experiences"
	experienceEventCollection = "experiences_events"
	collaboratorCollection    = "collaborators"
	commentCollection         = "comments"
	eventCollection           = "events"
)

type QueryModifier struct {
	SortBy  string
	SortAsc bool
	Skip    int
	Limit   int
}

func (m QueryModifier) modify(q *mgo.Query) *mgo.Query {
	if m.SortBy != "" {
		if !m.SortAsc {
			m.SortBy = "-" + m.SortBy
		}

		q = q.Sort(m.SortBy)
	}

	if m.Skip > 0 {
		q = q.Skip(m.Skip)
	}

	if m.Limit > 0 {
		q = q.Limit(m.Limit)
	}

	return q
}

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

	eventCollection := &EventCollection{db.C(eventCollection)}

	return &Database{
		session:         session,
		EventCollection: eventCollection,
		ExperienceCollection: &ExperienceCollection{
			experiences:       db.C(experienceCollection),
			experiencesEvents: db.C(experienceEventCollection),
			collaborators:     db.C(collaboratorCollection),
			comments:          db.C(commentCollection),
		},
	}, nil
}

func (db *Database) Close() {
	db.session.Close()
}

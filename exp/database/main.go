package database

import (
	"github.com/globalsign/mgo"
)

const (
	experienceTransactionCollection = "experiences.transactions"
	experienceCollection            = "experiences"
	experienceActivityCollection    = "experiences_activities"
	commentCollection               = "comments"
	activityCollection              = "activities"
)

type QueryModifier struct {
	SortBy  string `form:"sort_by"`
	SortAsc bool   `form:"sort_asc"`
	Skip    int    `form:"skip" binding:"omitempty,gt=0"`
	Limit   int    `form:"limit" binding:"omitempty,gt=0"`
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
	*ActivityCollection
	*ExperienceCollection
}

func NewDatabase(url string, dbName string) (*Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	db := session.DB(dbName)

	actCollection := &ActivityCollection{db.C(activityCollection)}

	return &Database{
		session:            session,
		ActivityCollection: actCollection,
		ExperienceCollection: &ExperienceCollection{
			transactions:          db.C(experienceTransactionCollection),
			experiences:           db.C(experienceCollection),
			comments:              db.C(commentCollection),
			experiencesActivities: db.C(experienceActivityCollection),
			activities:            actCollection,
		},
	}, nil
}

func (db *Database) Close() {
	db.session.Close()
}

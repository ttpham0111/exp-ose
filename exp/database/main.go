package database

import (
	"github.com/globalsign/mgo"
)

const (
	experienceTransactionCollection = "experiences.transactions"
	experienceCollection            = "experiences"
	ratingCollection                = "ratings"
	commentCollection               = "comments"
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
	*ExperienceCollection
}

func NewDatabase(url string, dbName string) (*Database, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	db := session.DB(dbName)

	return &Database{
		session: session,
		ExperienceCollection: &ExperienceCollection{
			transactions: db.C(experienceTransactionCollection),
			experiences:  db.C(experienceCollection),
			ratings:      db.C(ratingCollection),
			comments:     db.C(commentCollection),
		},
	}, nil
}

func (db *Database) Close() {
	db.session.Close()
}

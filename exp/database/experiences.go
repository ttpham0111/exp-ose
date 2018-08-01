package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const experienceType = "Experience"

type ExperienceCollectionReader interface {
	Find(ExperienceQuery, QueryModifier) ([]*Experience, error)
	FindId(string) (*Experience, error)
	GetEventsById(string) ([]*ExperienceEvent, error)
	GetCollaboratorsById(string) ([]UserId, error)
	GetCommentsById(string, QueryModifier) ([]*Comment, error)
}

type ExperienceCollection struct {
	experiences       *mgo.Collection
	experiencesEvents *mgo.Collection
	collaborators     *mgo.Collection
	comments          *mgo.Collection
}

type ExperienceQuery struct {
	Owner    string
	IsPublic bool
	Name     string
	Tags     []string
}

func (c *ExperienceCollection) Find(eq ExperienceQuery, m QueryModifier) ([]*Experience, error) {
	var fq = bson.M{
		"is_public": eq.IsPublic,
	}

	if eq.Owner != "" {
		fq["owner"] = eq.Owner
	}

	if eq.Name != "" {
		fq["name"] = eq.Name
	}

	if len(eq.Tags) == 1 {
		fq["tags"] = eq.Tags[0]
	} else if len(eq.Tags) > 1 {
		fq["tags"] = bson.M{"$all": eq.Tags}
	}

	if m.Limit == 0 {
		m.Limit = 100
	}

	experiences := make([]*Experience, 0)
	err := m.modify(c.experiences.Find(fq)).All(&experiences)
	return experiences, err
}

func (c *ExperienceCollection) FindId(experienceId string) (experience *Experience, err error) {
	if !bson.IsObjectIdHex(experienceId) {
		return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	err = c.experiences.FindId(bson.ObjectIdHex(experienceId)).One(experience)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
		}
	}

	return experience, err
}

func (c *ExperienceCollection) GetEventsById(experienceId string) ([]*ExperienceEvent, error) {
	if !bson.IsObjectIdHex(experienceId) {
		return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"experience_id": bson.ObjectIdHex(experienceId)},
		},
		{
			"$lookup": bson.M{
				"from":         eventCollection,
				"localField":   "event_id",
				"foreignField": "_id",
				"as":           "events",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"$mergeObjects": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$events", 0}}, "$$ROOT",
					},
				},
			},
		},
		{
			"$project": bson.M{"events": 0},
		},
	}
	events := make([]*ExperienceEvent, 0)
	err := c.experiencesEvents.Pipe(pipeline).All(&events)
	return events, err
}

func (c *ExperienceCollection) GetCollaboratorsById(experienceId string) ([]UserId, error) {
	if !bson.IsObjectIdHex(experienceId) {
		return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	collaborators := make([]UserId, 0)
	err := c.collaborators.Find(bson.M{"experience_id": bson.ObjectIdHex(experienceId)}).All(&collaborators)
	return collaborators, err
}

func (c *ExperienceCollection) GetCommentsById(experienceId string, m QueryModifier) ([]*Comment, error) {
	if !bson.IsObjectIdHex(experienceId) {
		return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	comments := make([]*Comment, 0)

	if m.SortBy == "" {
		m.SortBy = "created_at"
	}

	if m.Limit == 0 {
		m.Limit = 100
	}

	err := m.modify(c.comments.Find(bson.M{"experience_id": bson.ObjectIdHex(experienceId)})).All(&comments)
	return comments, err
}

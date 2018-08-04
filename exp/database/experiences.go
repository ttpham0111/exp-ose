package database

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo/txn"
)

const experienceType = "Experience"
const commentType = "Comment"
const expActivityType = "Experience Activity"

type ExperienceCollectionReadWriter interface {
	ExperienceCollectionReader
	ExperienceCollectionWriter
}

type ExperienceCollectionReader interface {
	Find(ExperienceQuery, QueryModifier) ([]*Experience, error)
	FindId(string) (*Experience, error)
	GetCommentsById(string, QueryModifier) ([]*Comment, error)
}

type ExperienceCollectionWriter interface {
	Create(*Experience) error
	Update(string, *Experience) error

	// AddRating(string, *Rating) (string, error)
	// RemoveRating(string, UserId) error

	AddComment(string, *Comment) error
	RemoveComment(string, string) error

	AddActivity(string, *ExperienceActivity) error
	RemoveActivity(string, string) error
}

type ExperienceCollection struct {
	transactions          *mgo.Collection
	experiences           *mgo.Collection
	comments              *mgo.Collection
	experiencesActivities *mgo.Collection
	activities            *ActivityCollection
}

type ExperienceQuery struct {
	Owner     string   `form:"owner"`
	IsPrivate bool     `form:"is_private"`
	Name      string   `form:"name"`
	Tags      []string `form:"tags" binding:"omitempty,gt=0"`
}

func (c *ExperienceCollection) Find(eq ExperienceQuery, m QueryModifier) ([]*Experience, error) {
	var fq = bson.M{
		"is_private": eq.IsPrivate,
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

func (c *ExperienceCollection) FindId(experienceId string) (*Experience, error) {
	if !bson.IsObjectIdHex(experienceId) {
		return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	var experience Experience
	err := c.experiences.FindId(bson.ObjectIdHex(experienceId)).One(&experience)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
		}
	}

	return &experience, nil
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

func (c *ExperienceCollection) GetActivitiesById(experienceId string) ([]*ExperienceActivity, error) {
	if !bson.IsObjectIdHex(experienceId) {
		return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"experience_id": bson.ObjectIdHex(experienceId)},
		},
		{
			"$lookup": bson.M{
				"localField":   "activity_id",
				"from":         activityCollection,
				"foreignField": "_id",
				"as":           "activities",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"$mergeObjects": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$activities", 0}}, "$$ROOT",
					},
				},
			},
		},
		{
			"$project": bson.M{"activities": 0},
		},
	}

	activities := make([]*ExperienceActivity, 0)
	err := c.experiencesActivities.Pipe(pipeline).All(&activities)
	return activities, err
}

func (c *ExperienceCollection) Create(experience *Experience) error {
	experience.Id = bson.NewObjectId()
	if experience.Tags == nil {
		experience.Tags = make([]string, 0)
	}
	if experience.Collaborators == nil {
		experience.Collaborators = make([]UserId, 0)
	}

	if err := c.experiences.Insert(experience); err != nil {
		return err
	}
	return nil
}

func (c *ExperienceCollection) Update(experienceId string, experience *Experience) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	experience.Id = bson.ObjectIdHex(experienceId)
	_, err := c.experiences.Upsert(nil, experience)
	return err
}

// TODO
// func (c *ExperienceCollection) AddRating(experienceId string, rating float32) (float32, error) {

// }

func (c *ExperienceCollection) AddComment(experienceId string, comment *Comment) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	comment.Id = bson.NewObjectId()
	comment.ExperienceId = bson.ObjectIdHex(experienceId)
	comment.CreatedAt = time.Now()

	runner := txn.NewRunner(c.transactions)
	ops := []txn.Op{
		{
			C:      commentCollection,
			Id:     comment.Id,
			Insert: comment,
		},
		{
			C:      experienceCollection,
			Id:     comment.ExperienceId,
			Update: bson.M{"$inc": bson.M{"num_comments": 1}},
		},
	}
	return runner.Run(ops, "", nil)
}

func (c *ExperienceCollection) RemoveComment(experienceId string, commentId string) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	if !bson.IsObjectIdHex(commentId) {
		return NoResultFound{ObjectType: commentType, ObjectId: commentId}
	}

	runner := txn.NewRunner(c.transactions)
	ops := []txn.Op{
		{
			C:      commentCollection,
			Id:     bson.ObjectIdHex(commentId),
			Remove: true,
		},
		{
			C:      experienceCollection,
			Id:     bson.ObjectIdHex(experienceId),
			Update: bson.M{"$inc": bson.M{"num_comments": -1}},
		},
	}
	return runner.Run(ops, "", nil)
}

func (c *ExperienceCollection) AddActivity(experienceId string, expActivity *ExperienceActivity) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	act := activity{activityBody: expActivity.activityBody}
	if err := c.activities.Create(&act); err != nil {
		return err
	}

	expActivity.experienceActivity.Id = bson.NewObjectId()
	expActivity.experienceActivity.ExperienceId = bson.ObjectIdHex(experienceId)
	expActivity.experienceActivity.ActivityId = act.Id

	runner := txn.NewRunner(c.transactions)
	ops := []txn.Op{
		{
			C:      experienceActivityCollection,
			Id:     expActivity.experienceActivity.Id,
			Insert: expActivity.experienceActivity,
		},
		{
			C:      experienceCollection,
			Id:     expActivity.experienceActivity.ExperienceId,
			Update: bson.M{"$inc": bson.M{"num_activities": 1}},
		},
	}
	return runner.Run(ops, "", nil)
}

func (c *ExperienceCollection) RemoveActivity(experienceId string, expActivityId string) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	if !bson.IsObjectIdHex(expActivityId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: expActivityId}
	}

	runner := txn.NewRunner(c.transactions)
	ops := []txn.Op{
		{
			C:      experienceActivityCollection,
			Id:     bson.ObjectIdHex(expActivityId),
			Remove: true,
		},
		{
			C:      experienceCollection,
			Id:     bson.ObjectIdHex(experienceId),
			Update: bson.M{"$inc": bson.M{"num_activities": -1}},
		},
	}
	return runner.Run(ops, "", nil)
}

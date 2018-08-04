package database

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const experienceType = "Experience"
const commentType = "Comment"

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

	AddComment(string, *Comment) error
	RemoveComment(string) error
}

type ExperienceCollection struct {
	transactions *mgo.Collection
	experiences  *mgo.Collection
	ratings      *mgo.Collection
	comments     *mgo.Collection
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
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
		}
	}

	return comments, err
}

func (c *ExperienceCollection) Create(experience *Experience) error {
	experience.Id = bson.NewObjectId()
	experience.initializeIfNil()
	experience.CreatedAt = time.Now()
	return c.experiences.Insert(experience)
}

func (c *ExperienceCollection) Update(experienceId string, experience *Experience) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	experience.Id = bson.ObjectIdHex(experienceId)
	experience.initializeIfNil()
	err := c.experiences.UpdateId(experience.Id, experience)
	if err != nil {
		if err == mgo.ErrNotFound {
			return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
		}
	}
	return err
}

func (c *ExperienceCollection) AddComment(experienceId string, comment *Comment) error {
	if !bson.IsObjectIdHex(experienceId) {
		return NoResultFound{ObjectType: experienceType, ObjectId: experienceId}
	}

	comment.Id = bson.NewObjectId()
	comment.ExperienceId = bson.ObjectIdHex(experienceId)
	comment.CreatedAt = time.Now()
	return c.comments.Insert(comment)
}

func (c *ExperienceCollection) RemoveComment(commentId string) error {
	if !bson.IsObjectIdHex(commentId) {
		return NoResultFound{ObjectType: commentType, ObjectId: commentId}
	}

	return c.comments.RemoveId(bson.ObjectIdHex(commentId))
}

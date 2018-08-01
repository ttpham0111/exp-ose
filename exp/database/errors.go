package database

type NoResultFound struct {
	ObjectType string
	ObjectId   string
}

func (err NoResultFound) Error() string {
	return err.ObjectType + " " + err.ObjectId + " not found"
}

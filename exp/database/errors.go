package database

type NoResultFound struct {
	ObjectType string
	ObjectId   string
}

func (err NoResultFound) Error() string {
	return err.ObjectType + " " + err.ObjectId + " not found"
}

type ValidationError struct {
	Field string
}

func (err ValidationError) Error() string {
	return "invalid value for " + err.Field
}

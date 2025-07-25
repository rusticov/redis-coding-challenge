package store

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrorKeyNotFound        Error = "key not found"
	ErrorNotAnInteger       Error = "not an integer"
	ErrorWrongOperationType Error = "wrong operation type"
)

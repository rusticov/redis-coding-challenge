package list

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrorOldValueIsNotList Error = "old value is not a list"
)

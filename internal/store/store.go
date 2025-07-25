package store

type Store interface {
	ReadString(key string) (string, error)
	ReadListRange(key string, fromIndex int, toIndex int) ([]string, error)
	Exists(key string) bool

	Write(key string, value string)
	Delete(key string) bool

	Increment(key string, incrementBy int64) (int64, error)
	LeftPush(key string, values []string) (int64, error)
	RightPush(key string, values []string) (int64, error)
}

type ExpiryOption string

const (
	ExpiryOptionNone                         ExpiryOption = ""
	ExpiryOptionExpirySeconds                ExpiryOption = "EX"
	ExpiryOptionExpiryMilliseconds           ExpiryOption = "PX"
	ExpiryOptionExpiryUnixTimeInSeconds      ExpiryOption = "EXAT"
	ExpiryOptionExpiryUnixTimeInMilliseconds ExpiryOption = "PXAT"
	ExpiryOptionExpiryKeepTTL                ExpiryOption = "KEEPTTL"
)

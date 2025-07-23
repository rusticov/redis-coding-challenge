package store

type Entry struct {
	data   any
	option ExpiryOption
	value  int64
}

func NewEntryWithExpiry(data any, option ExpiryOption, value int64) Entry {
	return Entry{data: data, option: option, value: value}
}

func NewEntry(data any) Entry {
	return Entry{data: data}
}

func (e Entry) Data() any {
	return e.data
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

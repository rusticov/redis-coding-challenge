package store

type ExpiryOption string

const (
	ExpiryOptionNone                         ExpiryOption = ""
	ExpiryOptionExpirySeconds                ExpiryOption = "EX"
	ExpiryOptionExpiryMilliseconds           ExpiryOption = "PX"
	ExpiryOptionExpiryUnixTimeInSeconds      ExpiryOption = "EXAT"
	ExpiryOptionExpiryUnixTimeInMilliseconds ExpiryOption = "PXAT"
	ExpiryOptionExpiryKeepTTL                ExpiryOption = "KEEPTTL"
)

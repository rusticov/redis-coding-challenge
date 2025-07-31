package command

import "redis-challenge/internal/store"

func ExpiryTimestamp(clock store.Clock, option store.ExpiryOption, expiry int64) (store.ExpiryOption, int64) {
	if option == store.ExpiryOptionNone || option == store.ExpiryOptionExpiryKeepTTL {
		return option, 0
	}

	var timestamp int64

	switch option {
	case store.ExpiryOptionExpiryUnixTimeInMilliseconds:
		timestamp = expiry
	case store.ExpiryOptionExpiryUnixTimeInSeconds:
		timestamp = expiry * 1000
	case store.ExpiryOptionExpiryMilliseconds:
		timestamp = clock() + expiry
	case store.ExpiryOptionExpirySeconds:
		timestamp = clock() + expiry*1000
	}

	return store.ExpiryOptionExpiryUnixTimeInMilliseconds, timestamp
}

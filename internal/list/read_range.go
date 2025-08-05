package list

func ReadRangeFromStoreList(storedData any, start, end int) (DoubleEndedList, bool) {
	storedList, ok := parseListFromStoredData(storedData)
	if !ok {
		return DoubleEndedList{}, false
	}

	return storedList.Filter(start, end), true
}

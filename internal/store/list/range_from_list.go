package list

func ReadRangeFromStoreList(storedData any, start, end int) ([]string, bool) {
	storedList, ok := parseListFromStoredData(storedData)
	if !ok {
		return nil, false
	}

	from := start
	if start < 0 {
		from = len(storedList) + start
	}
	if from < 0 {
		from = 0
	}

	to := end + 1
	if end < 0 {
		to = len(storedList) + end + 1
	}

	if storedList == nil {
		return nil, true
	}
	if to <= from {
		return nil, true
	}
	if to > len(storedList) {
		to = len(storedList)
	}
	return storedList[from:to], true
}

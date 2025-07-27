package list

func ReadRangeFromStoreList(storeValue any, start, end int) ([]string, bool) {
	allValues, ok := parseList(storeValue)
	if !ok {
		return nil, false
	}

	from := start
	if start < 0 {
		from = len(allValues) + start
	}
	if from < 0 {
		from = 0
	}

	to := end + 1
	if end < 0 {
		to = len(allValues) + end + 1
	}

	if allValues == nil {
		return nil, true
	}
	if to <= from {
		return nil, true
	}
	if to > len(allValues) {
		to = len(allValues)
	}
	return allValues[from:to], true
}

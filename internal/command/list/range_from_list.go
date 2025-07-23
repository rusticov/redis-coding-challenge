package list

func ReadRangeFromStoreList(storeValue any, start, end int) ([]string, error) {
	allValues, err := parseOldList(storeValue)
	if err != nil {
		return nil, err
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
		to = len(allValues) + end
	}

	if allValues == nil {
		return nil, nil
	}
	if to <= from {
		return nil, nil
	}
	if to > len(allValues) {
		to = len(allValues)
	}
	return allValues[from:to], nil
}

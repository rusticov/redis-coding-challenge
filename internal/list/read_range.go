package list

func ReadRangeFromStoreList(storedData any, start, end int) ([]string, bool) {
	storedList, ok := parseListFromStoredData(storedData)
	if !ok {
		return nil, false
	}

	length := storedList.Length()

	from := start
	if start < 0 {
		from = length + start
	}
	if from < 0 {
		from = 0
	}

	to := end + 1
	if end < 0 {
		to = length + end + 1
	}

	if length == 0 {
		return nil, true
	}
	if to <= from {
		return nil, true
	}
	if to > length {
		to = length
	}

	middleIndex := len(storedList.left)

	switch {
	case to <= middleIndex:
		left := make([]string, to-from)
		for i, x := range storedList.left[middleIndex-to : middleIndex-from] {
			left[to-from-i-1] = x
		}
		return left, true
	case from < middleIndex:
		result := make([]string, to-from)
		fromLeftCount := middleIndex - from

		for i, x := range storedList.left[0 : middleIndex-from] {
			result[fromLeftCount-i-1] = x
		}
		copy(result[fromLeftCount:], storedList.right)

		return result, true
	default:
		return storedList.right[from-middleIndex : to-middleIndex], true
	}
}

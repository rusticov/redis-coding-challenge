package list

func LeftPushToOldList(newValues []string, oldList any) ([]string, bool) {
	oldValues, ok := parseList(oldList)
	if !ok {
		return nil, false
	}

	newList := make([]string, len(newValues)+len(oldValues))
	copy(newList[len(newValues):], oldValues)

	newValueEndIndex := len(newValues) - 1

	for i, value := range newValues {
		newList[newValueEndIndex-i] = value
	}
	return newList, true
}

func RightPushToOldList(newValues []string, oldList any) ([]string, bool) {
	oldValues, ok := parseList(oldList)
	if !ok {
		return nil, false
	}

	newList := make([]string, len(newValues)+len(oldValues))
	copy(newList, oldValues)

	newValueStartIndex := len(oldValues)

	for i, value := range newValues {
		newList[newValueStartIndex+i] = value
	}
	return newList, true
}

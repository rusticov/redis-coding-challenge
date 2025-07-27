package list

func LeftPushToOldList(newValues []string, storedData any) ([]string, bool) {
	storedList, ok := parseListFromStoredData(storedData)
	if !ok {
		return nil, false
	}

	newList := make([]string, len(newValues)+len(storedList))
	copy(newList[len(newValues):], storedList)

	newValueEndIndex := len(newValues) - 1

	for i, value := range newValues {
		newList[newValueEndIndex-i] = value
	}
	return newList, true
}

func RightPushToOldList(newValues []string, storedData any) ([]string, bool) {
	storedList, ok := parseListFromStoredData(storedData)
	if !ok {
		return nil, false
	}

	newList := make([]string, len(newValues)+len(storedList))
	copy(newList, storedList)

	newValueStartIndex := len(storedList)

	for i, value := range newValues {
		newList[newValueStartIndex+i] = value
	}
	return newList, true
}

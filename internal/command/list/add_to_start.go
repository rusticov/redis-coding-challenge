package list

func LeftPushToOldList(newValues []string, oldList any) ([]string, error) {
	oldValues, err := parseOldList(oldList)
	if err != nil {
		return nil, err
	}

	newList := make([]string, len(newValues)+len(oldValues))
	copy(newList[len(newValues):], oldValues)

	newValueEndIndex := len(newValues) - 1

	for i, value := range newValues {
		newList[newValueEndIndex-i] = value
	}
	return newList, nil
}

func parseOldList(oldList any) ([]string, error) {
	if oldList == nil {
		return nil, nil
	}
	if stringList, ok := oldList.([]string); ok {
		return stringList, nil
	}
	return nil, ErrorOldValueIsNotList
}

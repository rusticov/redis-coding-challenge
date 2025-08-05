package list

func LeftPushToOldList(newValues []string, storedData any) (DoubleEndedList, bool) {
	list, ok := parseListFromStoredData(storedData)
	if !ok {
		return DoubleEndedList{}, false
	}

	list.Left = append(list.Left, newValues...)
	return list, true
}

func RightPushToOldList(newValues []string, storedData any) (DoubleEndedList, bool) {
	list, ok := parseListFromStoredData(storedData)
	if !ok {
		return DoubleEndedList{}, false
	}

	list.Right = append(list.Right, newValues...)
	return list, true
}

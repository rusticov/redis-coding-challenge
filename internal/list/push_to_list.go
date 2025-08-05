package list

func LeftPush(newValues []string, storedData any) (DoubleEndedList, bool) {
	list, ok := parseListFromStoredData(storedData)
	if !ok {
		return DoubleEndedList{}, false
	}

	list.left = append(list.left, newValues...)
	return list, true
}

func RightPush(newValues []string, storedData any) (DoubleEndedList, bool) {
	list, ok := parseListFromStoredData(storedData)
	if !ok {
		return DoubleEndedList{}, false
	}

	list.right = append(list.right, newValues...)
	return list, true
}

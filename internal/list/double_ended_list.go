package list

type DoubleEndedList struct {
	left  []string
	right []string
}

func (l DoubleEndedList) Length() int {
	return len(l.left) + len(l.right)
}

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

func parseListFromStoredData(storedData any) (DoubleEndedList, bool) {
	if storedData == nil {
		return DoubleEndedList{}, true
	}
	if stringList, ok := storedData.(DoubleEndedList); ok {
		return stringList, true
	}
	return DoubleEndedList{}, false
}

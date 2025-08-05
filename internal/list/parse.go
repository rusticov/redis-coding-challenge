package list

type DoubleEndedList struct {
	left  []string
	right []string
}

func (l DoubleEndedList) Length() int {
	return len(l.left) + len(l.right)
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

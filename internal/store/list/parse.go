package list

type DoubleEndedList struct {
	Left  []string
	Right []string
}

func (l DoubleEndedList) Length() int {
	return len(l.Left) + len(l.Right)
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

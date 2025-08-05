package list

import "iter"

type DoubleEndedList struct {
	left  []string
	right []string
}

func (l DoubleEndedList) Length() int {
	return len(l.left) + len(l.right)
}

func (l DoubleEndedList) Filter(start, end int) DoubleEndedList {
	length := l.Length()

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

	if length == 0 || to <= from {
		return DoubleEndedList{}
	}
	if to > length {
		to = length
	}

	middleIndex := len(l.left)

	switch {
	case to <= middleIndex:
		return DoubleEndedList{left: l.left[middleIndex-to : middleIndex-from]}
	case from < middleIndex:
		return DoubleEndedList{}
	default:
		return DoubleEndedList{right: l.right[from-middleIndex : to-middleIndex]}
	}
}

func (l DoubleEndedList) Range() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		i := 0

		for leftIndex := len(l.left) - 1; leftIndex >= 0; leftIndex-- {
			if !yield(i, l.left[leftIndex]) {
				return
			}
			i++
		}

		for _, x := range l.right {
			if !yield(i, x) {
				return
			}
			i++
		}
	}
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

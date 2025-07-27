package list

func parseList(oldList any) ([]string, bool) {
	if oldList == nil {
		return nil, true
	}
	if stringList, ok := oldList.([]string); ok {
		return stringList, true
	}
	return nil, false
}

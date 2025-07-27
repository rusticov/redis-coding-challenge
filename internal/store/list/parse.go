package list

func parseListFromStoredData(storedData any) ([]string, bool) {
	if storedData == nil {
		return nil, true
	}
	if stringList, ok := storedData.([]string); ok {
		return stringList, true
	}
	return nil, false
}

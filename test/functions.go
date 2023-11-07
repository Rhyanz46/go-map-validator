package test

func getAllKeys(data map[string]interface{}) (allKeysInMap []string) {
	for key, _ := range data {
		allKeysInMap = append(allKeysInMap, key)
	}
	return
}

func isDataInList[T comparable](key T, data []T) (result bool) {
	for _, val := range data {
		if val == key {
			return true
		}
	}
	return
}

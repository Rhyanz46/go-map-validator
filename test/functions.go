package test

import (
	"regexp"
	"strings"
)

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

func trimAndClean(data string) string {
	re := regexp.MustCompile(`[\t\n\r]+`)
	cleaned := re.ReplaceAllString(data, " ")
	trimmed := strings.TrimSpace(cleaned)
	words := strings.Fields(trimmed)
	return strings.Join(words, " ")
}

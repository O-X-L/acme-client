package u

import (
	"regexp"
	"slices"
)

func RegexMatch(p string, d string) bool {
	m, _ := regexp.MatchString(p, d)
	return m
}

func PtrBool(b bool) *bool {
	return &b
}

func RemoveDuplicates(input []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string

	for _, str := range input {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = true
			result = append(result, str)
		}
	}
	slices.Sort(result)
	return result
}

package u

import (
	"fmt"
	"slices"
	"strings"
)

func SortDomainsExceptFirst(input []string) {
	if len(input) < 2 {
		return
	}
	slices.Sort(input[1:])
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
	SortDomainsExceptFirst(result)
	return result
}

func BuildDiffMap(before, after []string) map[string][]string {
	diff := make(map[string][]string)
	beforeMap := make(map[string]bool)
	afterMap := make(map[string]bool)

	for _, s := range before {
		beforeMap[s] = true
	}
	for _, s := range after {
		afterMap[s] = true
	}

	for _, s := range before {
		if !afterMap[s] {
			diff["-"] = append(diff["-"], s)
		}
	}

	for _, s := range after {
		if !beforeMap[s] {
			diff["+"] = append(diff["+"], s)
		}
	}

	slices.Sort(diff["-"])
	slices.Sort(diff["+"])
	return diff
}

func BuildDiffString(diff map[string][]string) string {
	out := []string{}
	if diffAdd, exists := diff["+"]; exists && len(diff["+"]) > 0 {
		out = append(out, fmt.Sprintf("+[%s]", strings.Join(diffAdd, ", ")))
	}
	if diffRm, exists := diff["-"]; exists && len(diff["-"]) > 0 {
		out = append(out, fmt.Sprintf("-[%s]", strings.Join(diffRm, ", ")))
	}
	return strings.Join(out, " ")
}

func BuildBatches(input []string, size int) [][]string {
	out := [][]string{}
	count := len(input)
	for i := 0; i < count; i += size {
		end := min(i+size, count)
		out = append(out, input[i:end])
	}
	return out
}

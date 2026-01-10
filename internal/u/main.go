package u

import (
	"fmt"
	"os/user"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func RegexMatch(p string, d string) bool {
	m, _ := regexp.MatchString(p, d)
	return m
}

func PtrBool(b bool) *bool {
	return &b
}

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

func GetGroupID(grp string) (int, error) {
	if gid, err := strconv.Atoi(grp); err == nil {
		return gid, nil
	}

	if userGroup, err := user.LookupGroup(grp); err == nil {
		gid, _ := strconv.Atoi(userGroup.Gid)
		return gid, nil
	}

	return -1, fmt.Errorf("group '%s' is neither an existing group nor a GID", grp)
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
	for i := 0; i < len(input); i += size {
		end := i + size
		if end > len(input) {
			end = len(input)
		}
		out = append(out, input[i:end])
	}
	return out
}

package u

import (
	"fmt"
	"os/user"
	"regexp"
	"slices"
	"strconv"
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

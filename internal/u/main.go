package u

import (
	"fmt"
	"os/user"
	"regexp"
	"strconv"
)

func RegexMatch(p string, d string) bool {
	m, _ := regexp.MatchString(p, d)
	return m
}

func PtrBool(b bool) *bool {
	return &b
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

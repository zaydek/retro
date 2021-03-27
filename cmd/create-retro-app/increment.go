package create_retro_app

import (
	"regexp"
	"strconv"
)

var incrementRe = regexp.MustCompile(`(\d+)$`)

func increment(str string) string {
	if str == "" {
		return ""
	}
	matches := incrementRe.FindStringSubmatch(str)
	if matches == nil {
		return str + "2"
	}
	n, _ := strconv.Atoi(matches[1])
	return str[:len(str)-len(matches[1])] + strconv.Itoa(n+1)
}

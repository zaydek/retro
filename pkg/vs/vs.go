package vs

import (
	"errors"
	"regexp"
	"strings"
)

var ParseError = errors.New("versions: cannot parse contents")

type Version map[string]string

// https://regex101.com/r/0L0wqz/1
var rowRegex = regexp.MustCompile(`^\| +([^ ]+) +\| +([^ ]+) +\|$`)

func Parse(contents string) (Version, error) {
	vs := map[string]string{}

	arr := strings.Split(strings.TrimSpace(contents), "\n")
	if len(arr) < 2 {
		return nil, ParseError
	}
	subarr := arr[1 : len(arr)-1]
	for _, v := range subarr {
		matches := rowRegex.FindAllStringSubmatch(v, -1)
		if matches == nil {
			return nil, ParseError
		}
		pkg, v := matches[0][1], matches[0][2]
		vs[pkg] = v
	}
	return vs, nil
}

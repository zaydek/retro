package pretty

import (
	"encoding/json"
	"strings"
)

// PoorManJSON JSON encodes a value as:
//
// - "null"
// - "[foo, bar, baz]"
// - "{ "foo": "a", "bar": "b", "baz": "c" }"
//
func PoorManJSON(v interface{}) string {
	bstr, _ := json.MarshalIndent(v, "", " ")
	str := string(bstr)
	switch str[len(str)-1] {
	case ']':
		repl := strings.ReplaceAll(str, "\n", "")[2:] // Remove "[\n"
		return "[" + repl
	case '}':
		repl := strings.ReplaceAll(str, "\n", "")[2:] // Remove "{\n"
		return "{ " + repl[:len(repl)-1] + " }"
	}
	return str
}

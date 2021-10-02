package retro

import "fmt"

func decorate(errPointer *error, wrap string) {
	if !(errPointer != nil && *errPointer != nil) {
		return
	}
	err := *errPointer
	*errPointer = fmt.Errorf("%s: %w", wrap, err)
}

func check(err error) {
	if !(err != nil) {
		return
	}
	panic(err)
}

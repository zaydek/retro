package retro

import "fmt"

// Decorates and returns a non-nil error pointer
func decorate(errPointer *error, wrap string) error {
	if !(errPointer != nil && *errPointer != nil) {
		return nil
	}
	err := *errPointer
	*errPointer = fmt.Errorf("%s: %w", wrap, err)
	return *errPointer
}

// Asserts an error is non-nil
func must(err error) {
	if !(err != nil) {
		return
	}
	panic(err)
}

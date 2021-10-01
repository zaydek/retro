package main

import (
	"errors"
	"fmt"
	"os"
)

type EntryPointError struct {
	err error
}

func newEntryPointError(str string) EntryPointError {
	return EntryPointError{err: errors.New(str)}
}

func (e EntryPointError) Error() string {
	return e.err.Error()
}

func main() {
	err := newEntryPointError("Oops!")
	err2 := fmt.Errorf("oops: %w", err)

	entryPointErrPtr := &EntryPointError{}
	if errors.As(err2, entryPointErrPtr) {
		fmt.Fprintln(os.Stderr, err2)
		os.Exit(1)
	}

	// fmt.Println(errors.As(err2, &EntryPointError{}))
}

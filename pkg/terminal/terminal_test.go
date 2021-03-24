package terminal

import (
	"bytes"
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func TestRevert(t *testing.T) {
	var buf bytes.Buffer
	if _, err := Revert(&buf); err != nil {
		t.Fatal(err)
	}
	expect.DeepEqual(t, buf.String(), "\x1b[0m")
}

func TesetDeferRevert(t *testing.T) {
	var buf bytes.Buffer
	defer func() {
		expect.DeepEqual(t, buf.String(), "\x1b[0m")
	}()
	defer Revert(&buf)
}

func TestBold(t *testing.T) {
	expect.DeepEqual(t, Bold(), "")
	expect.DeepEqual(t, Boldf(""), "")
	expect.DeepEqual(t, Bold("Hello, world!"), "\x1b[1mHello, world!\x1b[0m")
	expect.DeepEqual(t, Boldf("%s", "Hello, world!"), "\x1b[1mHello, world!\x1b[0m")
}

func TestRed(t *testing.T) {
	expect.DeepEqual(t, Red(), "")
	expect.DeepEqual(t, Redf(""), "")
	expect.DeepEqual(t, Red("Hello, world!"), "\x1b[31mHello, world!\x1b[0m")
	expect.DeepEqual(t, Redf("%s", "Hello, world!"), "\x1b[31mHello, world!\x1b[0m")
}

func TestBoldRed(t *testing.T) {
	expect.DeepEqual(t, BoldRed(), "")
	expect.DeepEqual(t, BoldRedf(""), "")
	expect.DeepEqual(t, BoldRed("Hello, world!"), "\x1b[1m\x1b[31mHello, world!\x1b[0m")
	expect.DeepEqual(t, BoldRedf("%s", "Hello, world!"), "\x1b[1m\x1b[31mHello, world!\x1b[0m")
}

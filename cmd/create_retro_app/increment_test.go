package create_retro_app

import (
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func TestIncrement(t *testing.T) {
	expect.DeepEqual(t, increment(""), "")
	expect.DeepEqual(t, increment("hello-world"), "hello-world2")
	expect.DeepEqual(t, increment("hello-world0"), "hello-world1")
	expect.DeepEqual(t, increment("hello-world1"), "hello-world2")
	expect.DeepEqual(t, increment("hello-world2"), "hello-world3")
	expect.DeepEqual(t, increment("hello-world99"), "hello-world100")
	expect.DeepEqual(t, increment("hello-world99.99"), "hello-world99.100")
}

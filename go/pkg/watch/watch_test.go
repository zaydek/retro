package watch

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/zaydek/retro/go/pkg/expect"
)

func must(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatal(err)
}

func TestDirectory(t *testing.T) {
	var count int

	dir, err := ioutil.TempDir(".", "tmp_")
	must(t, err)
	defer os.RemoveAll(dir)

	ch := Directory(dir, 10*time.Millisecond)
	go func() {
		for range ch {
			count++
		}
	}()

	must(t, os.WriteFile(path.Join(dir, "a"), []byte("Hello, world!\n"), 0644))
	time.Sleep(10 * time.Millisecond)
	must(t, os.WriteFile(path.Join(dir, "b"), []byte("Hello, world!\n"), 0644))
	time.Sleep(10 * time.Millisecond)
	must(t, os.WriteFile(path.Join(dir, "c"), []byte("Hello, world!\n"), 0644))
	time.Sleep(10 * time.Millisecond)

	expect.DeepEqual(t, count, 3)
}

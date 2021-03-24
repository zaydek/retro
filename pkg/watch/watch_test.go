package watch

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/zaydek/retro/pkg/expect"
)

type Test struct{ got, want int }

func check(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatalf("check: %s", err)
}

func TestDirectory(t *testing.T) {
	var count int

	dir, err := ioutil.TempDir(".", "tmp_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ch := Directory(dir, 10*time.Millisecond)
	go func() {
		for range ch {
			count++
		}
	}()

	check(t, ioutil.WriteFile(path.Join(dir, "a"), []byte("Hello, world!\n"), 0644))
	time.Sleep(10 * time.Millisecond)
	check(t, ioutil.WriteFile(path.Join(dir, "b"), []byte("Hello, world!\n"), 0644))
	time.Sleep(10 * time.Millisecond)
	check(t, ioutil.WriteFile(path.Join(dir, "c"), []byte("Hello, world!\n"), 0644))
	time.Sleep(10 * time.Millisecond)

	expect.DeepEqual(t, count, 3)
}

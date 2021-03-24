package watch

import (
	"os"
	"path/filepath"
	"time"
)

type WatchResult struct{ Err error }

// Directory creates a new watcher for directory dir.
func Directory(dir string, poll time.Duration) <-chan WatchResult {
	var (
		ch       = make(chan WatchResult)
		mtimeMap = map[string]time.Time{}
	)

	go func() {
		defer close(ch)

		// Use time.NewTicker(poll) not time.Tick(poll); time.NewTicker(poll)
		// starts eagerly (see https://stackoverflow.com/a/47448177)
		ticker := time.NewTicker(poll)
		defer ticker.Stop()
		for ; true; <-ticker.C {
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if prev, ok := mtimeMap[path]; !ok {
					mtimeMap[path] = info.ModTime()
				} else {
					if next := info.ModTime(); prev != next {
						mtimeMap[path] = next
						ch <- WatchResult{nil}
					}
				}
				return nil
			})
			if err != nil {
				ch <- WatchResult{err}
			}
		}
	}()

	time.Sleep(time.Millisecond)
	return ch
}

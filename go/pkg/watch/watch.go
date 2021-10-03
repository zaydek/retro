package watch

import (
	"os"
	"path/filepath"
	"time"
)

type WatchResult struct {
	Err error
}

func Directory(dir string, poll time.Duration) <-chan WatchResult {
	var (
		ch         = make(chan WatchResult)
		modTimeMap = map[string]time.Time{}
	)

	go func() {
		defer close(ch)

		//https://stackoverflow.com/a/47448177
		ticker := time.NewTicker(poll)
		defer ticker.Stop()
		for ; true; <-ticker.C {
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if prev, ok := modTimeMap[path]; !ok {
					modTimeMap[path] = info.ModTime()
				} else {
					if next := info.ModTime(); prev != next {
						modTimeMap[path] = next
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

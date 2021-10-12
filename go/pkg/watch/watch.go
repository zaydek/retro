package watch

import (
	"io/fs"
	"path/filepath"
	"time"
)

type WatchEvent struct {
	Err error
}

func Directory(dir string, poll time.Duration) <-chan WatchEvent {
	var (
		ch         = make(chan WatchEvent)
		modTimeMap = map[string]time.Time{}
	)

	go func() {
		defer close(ch)

		// https://stackoverflow.com/a/47448177
		ticker := time.NewTicker(poll)
		defer ticker.Stop()
		for ; true; <-ticker.C {
			err := filepath.WalkDir(dir, func(root string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				info, _ := d.Info()
				if prev, ok := modTimeMap[root]; !ok {
					modTimeMap[root] = info.ModTime()
				} else {
					if next := info.ModTime(); prev != next {
						modTimeMap[root] = next
						ch <- WatchEvent{nil}
					}
				}
				return nil
			})
			if err != nil {
				ch <- WatchEvent{err}
			}
		}
	}()

	time.Sleep(time.Millisecond)
	return ch
}

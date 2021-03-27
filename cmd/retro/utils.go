package retro

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type lsInfo struct {
	path string
	size int64
}

type lsInfos []lsInfo

var greedyExtRe = regexp.MustCompile(`(\.).*$`)

func greedyExt(path string) string {
	matches := greedyExtRe.FindAllString(path, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

func (a lsInfos) Len() int      { return len(a) }
func (a lsInfos) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// // Sort by ext
// func (a lsInfos) Less(i, j int) bool { return greedyExt(a[i].path) < greedyExt(a[j].path) }

// Sort by name
func (a lsInfos) Less(i, j int) bool { return a[i].path < a[j].path }

func ls(dir string) (lsInfos, error) {
	var ls lsInfos
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ls = append(ls, lsInfo{
			path: path,
			size: info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ls, nil
}

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format
func byteCount(b int64) string {
	const u = 1024

	if b < u {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(u), 0
	for n := b / u; n >= u; n /= u {
		div *= u
		exp++
	}
	return fmt.Sprintf("%.0f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

package fsUtils

import (
	"os"
	"path/filepath"
)

type lsInfo struct {
	Path string
	Size int64
}

type lsInfos []lsInfo

func (a lsInfos) Len() int      { return len(a) }
func (a lsInfos) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Sort by name
func (a lsInfos) Less(i, j int) bool { return a[i].Path < a[j].Path }

func List(dir string) (lsInfos, error) {
	var ls lsInfos
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ls = append(ls, lsInfo{
			Path: path,
			Size: info.Size(),
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return ls, nil
}

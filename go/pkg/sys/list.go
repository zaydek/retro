package sys

import (
	"io/fs"
	"path/filepath"
)

type FileType int

const (
	File      FileType = 0
	Directory FileType = 1
)

type lsInfo struct {
	Type FileType
	Path string
	Size int64
}

type ls []lsInfo

func (a ls) Len() int      { return len(a) }
func (a ls) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Sort by name
func (a ls) Less(i, j int) bool { return a[i].Path < a[j].Path }

func List(dir string) (ls, error) {
	var ls ls
	err := filepath.WalkDir(dir, func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		var typ FileType
		if !d.IsDir() {
			typ = File
		} else {
			typ = Directory
		}
		info, _ := d.Info()
		ls = append(ls, lsInfo{
			Type: typ,
			Path: root,
			Size: info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ls, nil
}

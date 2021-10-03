package unix

import (
	"io/fs"
	"path/filepath"
)

type fileKind int

const (
	kindFile      fileKind = 0
	kindDirectory fileKind = 1
)

type lsInfo struct {
	// Unexported
	kind fileKind

	// Exported
	Path string
	Size int64
}

func (l lsInfo) IsDir() bool {
	return l.kind == kindDirectory
}

// type ls []lsInfo
//
// func (a ls) Len() int           { return len(a) }
// func (a ls) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a ls) Less(i, j int) bool { return a[i].Path < a[j].Path }

func List(dir string) ([]lsInfo, error) {
	var ls []lsInfo
	err := filepath.WalkDir(dir, func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		var kind fileKind
		if d.IsDir() {
			kind = kindDirectory
		} else {
			kind = kindFile
		}
		info, _ := d.Info()
		ls = append(ls, lsInfo{
			kind: kind,
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

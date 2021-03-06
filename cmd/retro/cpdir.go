package retro

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type copyInfo struct {
	source string
	target string
}

// err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
// 	if err != nil {
// 		return err
// 	}
// 	if info.IsDir() {
// 		return nil
// 	}
// 	for _, exclude := range excludes {
// 		if path == exclude {
// 			return nil
// 		}
// 	}
// 	cpInfos = append(cpInfos, copyInfo{
// 		source: path,
// 		target: filepath.Join(dst, filepath.Base(path)),
// 	})
// 	return nil
// })

func cpdir(src, dst string, excludes []string) error {
	// Sweep for sources and targets
	var cpInfos []copyInfo
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		for _, exclude := range excludes {
			if path == exclude {
				return nil
			}
		}
		cpInfos = append(cpInfos, copyInfo{
			source: path,
			target: filepath.Join(dst, filepath.Base(path)),
		})
		return nil
	})
	if err != nil {
		return err
	}

	// Copy sources to targets
	for _, cpInfo := range cpInfos {
		if dir := filepath.Dir(cpInfo.target); dir != "." {
			if err := os.MkdirAll(dir, MODE_DIR); err != nil {
				return err
			}
		}
		source, err := os.Open(cpInfo.source)
		if err != nil {
			return err
		}
		target, err := os.Create(cpInfo.target)
		if err != nil {
			return err
		}
		if _, err := io.Copy(target, source); err != nil {
			return err
		}
		source.Close()
		target.Close()
	}
	return nil
}

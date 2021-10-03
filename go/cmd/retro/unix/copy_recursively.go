package unix

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

func CopyRecursively(source, target string, excludes []string) error {
	// Sweep for sources and targets
	var copyInfos []copyInfo
	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
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
		copyInfos = append(copyInfos, copyInfo{
			source: path,
			target: filepath.Join(target, filepath.Base(path)),
		})
		return nil
	})
	if err != nil {
		return err
	}
	// Copy sources to targets
	for _, copyInfo := range copyInfos {
		// FIXME: It's not obvious what this does
		if filename := filepath.Dir(copyInfo.target); filename != "." {
			if err := os.MkdirAll(filename, 0755); err != nil {
				return err
			}
		}
		sourceFile, err := os.Open(copyInfo.source)
		if err != nil {
			return err
		}
		targetFile, err := os.Create(copyInfo.target)
		if err != nil {
			return err
		}
		if _, err := io.Copy(targetFile, sourceFile); err != nil {
			return err
		}
		sourceFile.Close()
		targetFile.Close()
	}
	return nil
}

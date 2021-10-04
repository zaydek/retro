package unix

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type cpInfo struct {
	srcPath string
	dstPath string
}

func CopyRecursively(srcDir, dstDir string, excludes []string) error {
	var cp []cpInfo
	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
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
		cp = append(cp, cpInfo{
			srcPath: path,
			dstPath: filepath.Join(dstDir, filepath.Base(path)),
		})
		return nil
	})
	if err != nil {
		return err
	}
	for _, info := range cp {
		if filename := filepath.Dir(info.dstPath); filename != "." {
			if err := os.MkdirAll(filename, 0755); err != nil {
				return err
			}
		}
		src, err := os.Open(info.srcPath)
		if err != nil {
			return err
		}
		dst, err := os.Create(info.dstPath)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		src.Close()
		dst.Close()
	}
	return nil
}

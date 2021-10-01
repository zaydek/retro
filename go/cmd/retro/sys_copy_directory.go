package retro

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/zaydek/retro/cmd/perm"
)

// TODO: In theory this can and should be extracted to a separate package since
// it has nothing to do with Retro

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

// Copies a directory recursively
func copyDirectory(src, dst string, excludes []string) error {
	// TODO: Do we want to guard for non-directory sources and or destinations?

	// Sweep for sources and targets
	var copyInfos []copyInfo
	if err := filepath.WalkDir(src, func(path string, directoryEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if directoryEntry.IsDir() {
			return nil
		}
		for _, exclude := range excludes {
			if path == exclude {
				return nil
			}
		}
		copyInfos = append(copyInfos, copyInfo{
			source: path,
			target: filepath.Join(dst, filepath.Base(path)),
		})
		return nil
	}); err != nil {
		return err
	}

	// Copy sources to targets
	for _, copyInfo := range copyInfos {
		// FIXME: It's not obvious what this does
		if filename := filepath.Dir(copyInfo.target); filename != "." {
			if err := os.MkdirAll(filename, perm.BitsDirectory); err != nil {
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

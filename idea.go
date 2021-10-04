package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

var (
	//go:embed go/cmd/create_retro_app/static/*
	wwwFS embed.FS

	// wwwFS fs.FS
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

// func init() {
// 	wwwFS, _ = fs.Sub(wwwFS, "go/cmd/create_retro_app/static")
// }

func main() {
	var paths []string
	err := fs.WalkDir(wwwFS, ".", func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, root)
		}
		return nil
	})
	must(err)
	for _, path := range paths {
		fmt.Println(path)
	}

	f, err := wwwFS.Open(filepath.Join("go/cmd/create_retro_app/static", ".gitignore"))
	// b := make([]byte, 1024)
	var b []byte
	_, err = f.Read(b)
	must(err)
	fmt.Println(string(b))
	// must(err)
	// bstr, err := ioutil.ReadAll(f)
	// must(err)
	// fmt.Println("Here", string(bstr))

	// var bstr []byte
	// _, err = f.Read(bstr)
	// must(err)
	// fmt.Println("Here", string(bstr))

	// _, err := os.Stat("www")
	// fmt.Println(err)
	// fmt.Println(os.IsExist(err))
	// return
	//
	// wd, _ := os.Getwd()
	// fmt.Println(wd)
	// return
	//
	// var paths []string
	// err = fs.WalkDir(wwwFS, ".", func(root string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !d.IsDir() {
	// 		paths = append(paths, root)
	// 	}
	// 	return nil
	// })
	// must(err)
	//
	// for _, path := range paths {
	// 	rel, _ := filepath.Rel("go/cmd/create_retro_app/static", path)
	// 	fmt.Println(rel)
	// }
}

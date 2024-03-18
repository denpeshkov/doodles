package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
)

var dirsCount, filesCount = 0, 0

func main() {
	log.SetFlags(0)

	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	for _, root := range args {
		fmt.Println(root)
		dirsCount++
		if err := tree(root, ""); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("\n%d directories, %d files\n", dirsCount, filesCount)
}

func tree(path string, indent string) error {
	entrs, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read directory: %q", path)
	}
	// delete hidden entries
	entrs = slices.DeleteFunc(entrs, func(e fs.DirEntry) bool { return e.Name()[0] == '.' })

	for i, e := range entrs {
		prefix := "├── "
		add := "│   "
		// if last element
		if i == len(entrs)-1 {
			prefix = "└── "
			add = "    "
		}
		fmt.Printf("%s%s\n", indent+prefix, e.Name())

		if !e.IsDir() {
			filesCount++
			continue
		}
		dirsCount++
		if err := tree(filepath.Join(path, e.Name()), indent+add); err != nil {
			return err
		}
	}
	return nil
}

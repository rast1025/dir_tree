package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func dirTree(w io.Writer, path string, printFiles bool) error {
	return walk(w, path, "", printFiles)

}

func removeItem(s []os.FileInfo, index int) []os.FileInfo {
	return append(s[:index], s[index+1:]...)
}

func walk(w io.Writer, path string, prefix string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if !printFiles {
		for i:=0; i< len(files); i++ {
			if !files[i].IsDir() {
				files = removeItem(files, i)
			}
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	for i, el := range files {
		if !el.IsDir() && !printFiles {
			continue
		}
		var size string
		if !el.IsDir() {
			if el.Size() > 0 {
				size = " (" + strconv.FormatInt(el.Size(), 10) + "b)"
			} else {
				size = " (empty)"
			}
		}
		if i+1 == len(files) {
			fmt.Fprintf(w, "%v└───%v%v\n", prefix, files[i].Name(), size)
		} else {
			fmt.Fprintf(w, "%v├───%v%v\n", prefix, files[i].Name(), size)
		}
		if el.IsDir() {
			var newPrefix string
			if i+1 == len(files) {
				newPrefix = prefix + "\t"
			} else {
				newPrefix = prefix + "│\t"
			}
			err := walk(w, path+"/"+el.Name(), newPrefix, printFiles)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

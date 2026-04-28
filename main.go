package main

import (
	"fmt"
	"io"
	"os"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	return helper(out, path, printFiles, "")
}

func helper(out io.Writer, path string, printFiles bool, prefix string) error {
	var (
		last      string = "└───"
		emptyFile string = "empty"
		notLast   string = "├───"
	)

	files, err := filter(path, printFiles)
	if err != nil {
		return err
	}

	for i, f := range files {
		info, err := f.Info()
		if err != nil {
			return err
		}

		sizeInfo := info.Size()

		if sizeInfo != 0 {
			emptyFile = fmt.Sprintf("%db", sizeInfo)
		}

		if f.IsDir() {
			if i == len(files)-1 {
				fmt.Fprintf(out, prefix+last+f.Name()+"\n")
				helper(out, path+"/"+f.Name(), printFiles, prefix+"\t")
			} else {
				fmt.Fprintf(out, prefix+notLast+f.Name()+"\n")
				helper(out, path+"/"+f.Name(), printFiles, prefix+"│\t")
			}

		} else {
			if i == len(files)-1 {
				fmt.Fprintf(out, prefix+last+f.Name()+" (%s)\n", emptyFile)
			} else {
				fmt.Fprintf(out, prefix+notLast+f.Name()+" (%s)\n", emptyFile)
			}
		}
	}

	return nil
}

func filter(path string, printFiles bool) ([]os.DirEntry, error) {
	out := []os.DirEntry{}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		if file.IsDir() {
			out = append(out, file)
		} else if printFiles {
			out = append(out, file)
		}
	}

	return out, nil
}

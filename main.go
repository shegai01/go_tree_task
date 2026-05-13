package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var (
		magicNum3 = 3
		magicNum2 = 2
	)

	out := os.Stdout

	if len(os.Args) != magicNum2 && len(os.Args) != magicNum3 {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == magicNum3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return buildTree(out, path, "", printFiles)
}

func buildTree(out io.Writer, path, prefix string, printFiles bool) error { //
	const (
		last    string = "└───"
		notLast string = "├───"
	)

	files, err := getFiles(path, printFiles)
	if err != nil {
		return err
	}

	var (
		newPrefix string
		branch    string
	)

	for i, f := range files {
		if isLast(files, i) {
			newPrefix = prefix + "\t"
			branch = last
		} else {
			newPrefix = prefix + "│\t"
			branch = notLast
		}

		if f.IsDir() {
			fmt.Fprintf(out, prefix+branch+f.Name()+"\n")

			if err = buildTree(out, path+"/"+f.Name(), newPrefix, printFiles); err != nil {
				return err
			}

		} else {
			sizeInfo, errSize := getSize(f)
			if err != nil {
				return errSize
			}

			fmt.Fprintf(out, prefix+branch+f.Name()+" (%s)\n", sizeInfo)
		}
	}

	return nil
}

func getFiles(path string, printFiles bool) ([]os.DirEntry, error) {
	var res []os.DirEntry

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || printFiles {
			res = append(res, file)
		}
	}

	return res, err
}

func getSize(file os.DirEntry) (string, error) {
	emptyFile := "empty"

	info, err := file.Info()
	if info.Size() == 0 {
		return emptyFile, nil
	}

	return fmt.Sprintf("%db", info.Size()), err
}

func isLast(arr []os.DirEntry, i int) bool {
	return len(arr)-1 == i
}

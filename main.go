package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	out := os.Stdout

	if len(os.Args) != 2 && len(os.Args) != 3 {
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
	return buildTree(out, path, "", printFiles)
}

func buildTree(out io.Writer, path, prefix string, printFiles bool) error { //
	const (
		last    string = "└───"
		notLast string = "├───"
	)

	files, err := getFile(path, printFiles)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			fmt.Fprintf(out, prefix+last+f.Name()+"\n")

			if err = buildTree(out, path+"/"+f.Name(), prefix+"\t", printFiles); err != nil {
				return err
			}
			// } else {
			// 	fmt.Fprintf(out, prefix+notLast+f.Name()+"\n")
			// }
			// 	if err := buildTree(out, path+"/"+f.Name(), prefix+"│\t", printFiles); err != nil {
			// 		return err
			// 	}

		} else {
			sizeInfo, errSize := getSize(f)
			if err != nil {
				return errSize
			}

			// fmt.Fprintf(out, prefix+last+f.Name()+" (%s)\n", sizeInfo)

			fmt.Fprintf(out, prefix+notLast+f.Name()+" (%s)\n", sizeInfo)
		}

	}
	return nil
}

func getFile(path string, printFiles bool) ([]os.DirEntry, error) {
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

	return nil, err
}

func getSize(file os.DirEntry) (string, error) {
	emptyFile := "empty"

	info, err := file.Info()
	if info.Size() != 0 {
		emptyFile = fmt.Sprintf("%db", info.Size())
		return emptyFile, nil
	}

	return "emptyFile", err
}

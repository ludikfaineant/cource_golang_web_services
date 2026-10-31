package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Ancestors []bool

func (x *Ancestors) Pop() (bool, error) {
	if len(*x) == 0 {
		return false, fmt.Errorf("Lenght stack is 0")
	}
	item := (*x)[len(*x)-1]
	*x = (*x)[:len(*x)-1]
	return item, nil
}

func (x *Ancestors) Push(item bool) {
	*x = append(*x, item)
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

func dirTree(out io.Writer, path string, printFiles bool) error {
	ancestors := Ancestors{}
	return dirTreeWalk(out, path, printFiles, ancestors)
}

func dirTreeWalk(out io.Writer, path string, printFiles bool, ancestors Ancestors) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})
	for i, x := range dir {
		var lastFile bool
		if printFiles {
			lastFile = (i == len(dir)-1)
		} else {
			lastFile = (i == findLastDir(dir))
		}
		if x.IsDir() {
			drawTree(out, ancestors)
			if lastFile {
				ancestors.Push(true)
				fmt.Fprint(out, "└───", x.Name(), "\n")
			} else {
				ancestors.Push(false)
				fmt.Fprint(out, "├───", x.Name(), "\n")
			}
			err = dirTreeWalk(out, path+"/"+x.Name(), printFiles, ancestors)
			if err != nil {
				return err
			}
			_, err := ancestors.Pop()
			if err != nil {
				return err
			}
		} else {
			if printFiles {
				drawTree(out, ancestors)
				info, err := x.Info()
				if err != nil {
					return err
				}

				if lastFile {
					fmt.Fprint(out, "└───")
				} else {
					fmt.Fprint(out, "├───")
				}
				if info.Size() > 0 {
					fmt.Fprint(out, x.Name(), " (", info.Size(), "b)\n")
				} else {
					fmt.Fprint(out, x.Name(), " (empty)\n")
				}
			}

		}

	}
	return nil
}

func drawTree(out io.Writer, ancestors Ancestors) {
	for _, x := range ancestors {
		if x == true {
			fmt.Fprint(out, "	")
		} else {
			fmt.Fprint(out, "│	")
		}
	}
}

func findLastDir(dir []os.DirEntry) int {
	for i := len(dir) - 1; i > 0; i-- {
		if dir[i].IsDir() {
			return i
		}
	}
	return 0
}

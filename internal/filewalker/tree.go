package filewalker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DirTree struct {
	Name     string
	Path     string
	Children []*DirTree
}

func BuildDirTree(root string) (*DirTree, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	var dirTree *DirTree
	if !info.IsDir() {
		parts := strings.Split(info.Name(), ".")
		ext := parts[len(parts)-1]
		if ext == "vm" {
			dirTree = &DirTree{
				Name: info.Name(),
				Path: root,
			}
		}
	} else {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}
		dirTree = &DirTree{
			Name: info.Name(),
			Path: root,
		}
		var children []*DirTree

		for _, entry := range entries {
			childPath := filepath.Join(root, entry.Name())
			childNode, err := BuildDirTree(childPath)
			if err != nil {
				return nil, err
			}
			if childNode != nil {
				children = append(children, childNode)
			}
		}
		if len(children) != 0 {
			dirTree.Children = children
		}
	}
	return dirTree, nil
}

func PrintDirTree(tree *DirTree, level int) {
	trail := ""
	for i := 0; i < level; i++ {
		trail += "-"
	}
	fmt.Printf("%s%s\n", trail, tree.Name)
	for i := 0; i < len(tree.Children); i++ {
		PrintDirTree(tree.Children[i], level+1)
	}
}

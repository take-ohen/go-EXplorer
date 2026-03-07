package ui

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/lxn/walk"
)

// DirectoryNode represents a node in the directory tree.
type DirectoryNode struct {
	path     string
	name     string
	parent   *DirectoryNode
	children []*DirectoryNode
}

// Path returns the full path of the node.
func (d *DirectoryNode) Path() string {
	return d.path
}

// Text returns the display name of the node.
func (d *DirectoryNode) Text() string {
	return d.name
}

// Parent returns the parent node.
func (d *DirectoryNode) Parent() walk.TreeItem {
	if d.parent == nil {
		return nil
	}
	return d.parent
}

// ChildCount returns the number of child nodes.
func (d *DirectoryNode) ChildCount() int {
	if d.children == nil {
		children, err := d.loadChildren()
		if err != nil {
			// Access denied or other errors are ignored for tree view performance
		}
		d.children = children
	}
	return len(d.children)
}

// ChildAt returns the child node at the given index.
func (d *DirectoryNode) ChildAt(index int) walk.TreeItem {
	return d.children[index]
}

func (d *DirectoryNode) loadChildren() ([]*DirectoryNode, error) {
	var children []*DirectoryNode

	entries, err := os.ReadDir(d.path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			// Skip hidden directories for clarity
			if strings.HasPrefix(name, ".") {
				continue
			}
			child := &DirectoryNode{
				path:   filepath.Join(d.path, name),
				name:   name,
				parent: d,
			}
			children = append(children, child)
		}
	}

	return children, nil
}

// DirectoryTreeModel is the model for the directory TreeView.
type DirectoryTreeModel struct {
	walk.TreeModelBase
	roots []*DirectoryNode
}

// NewDirectoryTreeModel creates and initializes a new DirectoryTreeModel.
func NewDirectoryTreeModel() *DirectoryTreeModel {
	m := new(DirectoryTreeModel)
	m.populateRoots()
	return m
}

func (m *DirectoryTreeModel) populateRoots() {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return
	}
	getLogicalDrives, err := kernel32.FindProc("GetLogicalDrives")
	if err != nil {
		return
	}

	ret, _, _ := getLogicalDrives.Call()
	bitMap := uint32(ret)

	for i := 0; i < 26; i++ {
		if bitMap&(1<<uint(i)) != 0 {
			path := string('A'+rune(i)) + ":\\"
			m.roots = append(m.roots, &DirectoryNode{path: path, name: path})
		}
	}
}

// RootCount returns the number of root nodes.
func (m *DirectoryTreeModel) RootCount() int {
	return len(m.roots)
}

// RootAt returns the root node at the given index.
func (m *DirectoryTreeModel) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}

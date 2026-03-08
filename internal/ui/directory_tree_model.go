package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/lxn/walk"
)

// DirectoryNode represents a node in the directory tree.
type DirectoryNode struct {
	path     string
	name     string
	parent   *DirectoryNode
	children []*DirectoryNode
	scanned  bool
	scanning bool
	model    *DirectoryTreeModel
	mu       sync.Mutex
}

func NewDirectoryNode(path, name string, parent *DirectoryNode, model *DirectoryTreeModel) *DirectoryNode {
	return &DirectoryNode{
		path:   path,
		name:   name,
		parent: parent,
		model:  model,
		// [Model] children は初期状態では nil
		children: nil,
		scanned:  false,
	}
}

func (d *DirectoryNode) Path() string { return d.path }
func (d *DirectoryNode) Text() string { return d.name }
func (d *DirectoryNode) Parent() walk.TreeItem {
	if d.parent == nil {
		return nil
	}
	return d.parent
}

// [Model] ChildCount returns 1 if not scanned to show "+", otherwise real count.
func (d *DirectoryNode) ChildCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	count := 0
	if !d.scanned {
		count = 1
	} else {
		count = len(d.children)
	}
	fmt.Printf("[DEBUG] ChildCount for %s: scanned=%v, return=%d\n", d.path, d.scanned, count)
	return count
}

// [Model] ChildAt scans if needed and returns the real child.
func (d *DirectoryNode) ChildAt(index int) walk.TreeItem {
	d.mu.Lock()
	defer d.mu.Unlock()

	fmt.Printf("[DEBUG] ChildAt(%d) for %s: scanned=%v\n", index, d.path, d.scanned)

	if !d.scanned {
		if !d.scanning {
			fmt.Printf("[DEBUG] Triggering Scan for %s\n", d.path)
			d.scanning = true
			go d.Scan()
		}
		// スキャン中はプレースホルダーを返す（walkの仕様上、nilを返すと落ちるため）
		placeholder := NewDirectoryNode("", "Loading...", d, d.model)
		placeholder.scanned = true // プレースホルダーはスキャン対象外とする（無限ループ防止）
		fmt.Printf("[DEBUG] Returning Placeholder for %s\n", d.path)
		return placeholder
	}

	if index >= 0 && index < len(d.children) {
		child := d.children[index]
		fmt.Printf("[DEBUG] Returning Child for %s: %s (Addr: %p)\n", d.path, child.name, child)
		return child
	}
	fmt.Printf("[DEBUG] ChildAt index out of bounds for %s: %d (len=%d)\n", d.path, index, len(d.children))
	return nil
}

func (d *DirectoryNode) Image() interface{} { return nil }

// [Concurrency] Scan executes directory scanning safely.
func (d *DirectoryNode) Scan() {
	fmt.Printf("[DEBUG] Scan start: %s\n", d.path)
	d.mu.Lock()
	if d.scanned {
		fmt.Printf("[DEBUG] Scan skipped (already scanned): %s\n", d.path)
		d.mu.Unlock()
		return
	}
	d.mu.Unlock() // I/O前にロック解除

	// ディレクトリ読み込み
	entries, err := os.ReadDir(d.path)
	if err != nil {
		fmt.Printf("[DEBUG] ReadDir failed for %s: %v\n", d.path, err)
	} else {
		fmt.Printf("[DEBUG] ReadDir success for %s: %d entries\n", d.path, len(entries))
	}

	d.mu.Lock() // 結果反映のため再ロック
	// 他のスレッドが完了させていないか確認
	if d.scanned {
		d.mu.Unlock()
		return
	}

	var newChildren []*DirectoryNode
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				if strings.HasPrefix(name, ".") {
					continue
				}
				newChildren = append(newChildren, NewDirectoryNode(filepath.Join(d.path, name), name, d, d.model))
			}
		}
	}

	d.children = newChildren
	d.scanned = true
	d.scanning = false
	fmt.Printf("[DEBUG] Scan processed for %s: %d children created\n", d.path, len(newChildren))
	for i, child := range newChildren {
		if i < 3 {
			fmt.Printf("[DEBUG]   Child[%d]: %s (Addr: %p)\n", i, child.name, child)
		}
	}
	d.mu.Unlock()
	fmt.Printf("[DEBUG] Scan finished: %s\n", d.path)

	// [Concurrency] 完了後に PublishItemChanged を呼ぶ
	if d.model != nil && d.model.mw != nil {
		d.model.mw.Synchronize(func() {
			fmt.Printf("[DEBUG] Notifying view (PublishItemsReset): %s\n", d.path)
			d.model.PublishItemsReset(d)
		})
	} else {
		log.Printf("[WARN] MainWindow is nil. UI update skipped for %s", d.path)
	}
}

type DirectoryTreeModel struct {
	walk.TreeModelBase
	roots []*DirectoryNode
	mw    *walk.MainWindow
}

func NewDirectoryTreeModel() *DirectoryTreeModel {
	m := &DirectoryTreeModel{}
	return m
}

// SetMainWindow sets the main window and populates the roots.
func (m *DirectoryTreeModel) SetMainWindow(mw *walk.MainWindow) {
	m.mw = mw
	m.populateRoots()
	m.PublishItemsReset(nil)
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
			m.roots = append(m.roots, NewDirectoryNode(path, path, nil, m))
		}
	}
}

func (m *DirectoryTreeModel) RootCount() int                 { return len(m.roots) }
func (m *DirectoryTreeModel) RootAt(index int) walk.TreeItem { return m.roots[index] }

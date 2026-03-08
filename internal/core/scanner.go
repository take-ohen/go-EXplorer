package core

import (
	"os"
	"path/filepath"
	"sync"
)

// FileItem represents a file or directory found during scanning.
type FileItem struct {
	Path  string
	Name  string
	IsDir bool
	Size  int64
}

// Scanner provides high-speed directory scanning capabilities.
type Scanner struct {
	// sem limits the number of concurrent goroutines to prevent resource exhaustion.
	sem chan struct{}
}

// NewScanner creates a new Scanner instance with a concurrency limit.
func NewScanner() *Scanner {
	// Limit concurrency to 32 to avoid "too many open files" errors.
	return &Scanner{
		sem: make(chan struct{}, 32),
	}
}

// Scan starts scanning the directory at root.
// It returns a channel for FileItem objects and a channel for errors.
// Both channels are closed when the scan is complete.
func (s *Scanner) Scan(root string) (<-chan FileItem, <-chan error) {
	out := make(chan FileItem, 1000)
	errc := make(chan error, 1) // Buffered to avoid blocking on first error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.walk(root, out, errc, &wg)
	}()

	// Monitor routine to close channels when all workers are done.
	go func() {
		wg.Wait()
		close(out)
		close(errc)
	}()

	return out, errc
}

func (s *Scanner) walk(path string, out chan<- FileItem, errc chan<- error, wg *sync.WaitGroup) {
	// Acquire semaphore token
	s.sem <- struct{}{}
	defer func() {
		// Release semaphore token
		<-s.sem
	}()

	entries, err := os.ReadDir(path)
	if err != nil {
		errc <- err
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		info, err := entry.Info()
		if err != nil {
			errc <- err
			// Do not continue; try to process even if Info() fails.
		}

		var size int64
		if info != nil {
			size = info.Size()
		}

		out <- FileItem{
			Path:  fullPath,
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  size,
		}

		if entry.IsDir() {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				s.walk(p, out, errc, wg)
			}(fullPath)
		}
	}
}

// ListDir returns the list of files and directories in the specified path.
// It does not traverse subdirectories.
func ListDir(path string) ([]FileItem, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []FileItem
	for _, entry := range entries {
		info, _ := entry.Info() // Error ignored for list view performance
		var size int64
		if info != nil {
			size = info.Size()
		}
		items = append(items, FileItem{Path: filepath.Join(path, entry.Name()), Name: entry.Name(), IsDir: entry.IsDir(), Size: size})
	}
	return items, nil
}

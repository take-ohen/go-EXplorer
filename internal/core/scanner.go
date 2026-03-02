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
// It returns a channel that receives FileItem objects as they are found.
// The channel is closed when the scan is complete.
func (s *Scanner) Scan(root string) <-chan FileItem {
	out := make(chan FileItem, 1000) // Buffered channel for performance
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.walk(root, out, &wg)
	}()

	// Monitor routine to close the channel when all workers are done.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *Scanner) walk(path string, out chan<- FileItem, wg *sync.WaitGroup) {
	// Acquire semaphore token
	s.sem <- struct{}{}
	defer func() {
		// Release semaphore token
		<-s.sem
	}()

	entries, err := os.ReadDir(path)
	if err != nil {
		// Permission denied or other errors are ignored in this prototype.
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		// Get file info for size
		info, err := entry.Info()
		var size int64
		if err == nil {
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
				s.walk(p, out, wg)
			}(fullPath)
		}
	}
}

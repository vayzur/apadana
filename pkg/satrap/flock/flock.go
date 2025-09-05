package flock

import (
	"fmt"
	"os"
	"syscall"
)

type Flock struct {
	path string
	file *os.File
}

func NewFlock(path string) *Flock {
	return &Flock{
		path: path,
	}
}

func (f *Flock) TryLock() error {
	if f.file != nil {
		return fmt.Errorf("lock already acquired")
	}

	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("lock already held by another process")
		}
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	fmt.Fprintf(file, "%d", os.Getpid())

	// file stays accessible via fd but auto-cleans on process death
	syscall.Unlink(f.path)

	f.file = file
	return nil
}

func (f *Flock) Lock() error {
	return f.TryLock()
}

func (f *Flock) Unlock() error {
	if f.file == nil {
		return fmt.Errorf("lock not held")
	}

	// Just close the fd - file is already unlinked
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("failed to close lock file: %w", err)
	}

	f.file = nil
	return nil
}

func (f *Flock) IsLocked() bool {
	return f.file != nil
}

func (f *Flock) Path() string {
	return f.path
}

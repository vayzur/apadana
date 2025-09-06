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

	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open lock file: %w", err)
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		if err == syscall.EWOULDBLOCK {
			return fmt.Errorf("lock already held by another process")
		}
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	file.Truncate(0)
	file.Seek(0, 0)
	fmt.Fprintf(file, "%d", os.Getpid())
	file.Sync()

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

	// Closing the file automatically releases the flock
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

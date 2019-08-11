package paths

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
)

// DFB returns the dfb path
func DFB() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/.dfb", home)
}

// Exists checks if path exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatal(err)
		return false
	}
}

// SymlinkExists checks if a symbolink link exist at given path
func SymlinkExists(path string) bool {
	if _, err := os.Lstat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatal(err)
		return false
	}
}

// IsDir checks if file at given path is a directory
func IsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	return file.Mode().IsDir()
}

func IsSymlink(path string) bool {
	if !SymlinkExists(path) {
		return false
	}

	file, err := os.Lstat(path)
	if err != nil {
		log.Fatal(err)
	}

	_, ok := file.Sys().(*syscall.Stat_t)
	if !ok {
		err = errors.New("cannot convert stat value to syscall.Stat_t")
		log.Fatal(err)
	}

	// True if the file is a symlink.
	if file.Mode()&os.ModeSymlink != 0 {
		_, err := os.Readlink(path)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
	return false
}

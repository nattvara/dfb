package paths

import (
	"fmt"
	"log"
	"os"
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

package util

import (
	"log"
	"os"
	"path/filepath"
)

// FindProjectRoot ищет корень проекта по наличию папки data или static
func FindProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get cwd: %v", err)
	}
	dir := cwd
	for i := 0; i < 4; i++ {
		if fi, err := os.Stat(filepath.Join(dir, "data")); err == nil && fi.IsDir() {
			return dir
		}
		if fi, err := os.Stat(filepath.Join(dir, "static")); err == nil && fi.IsDir() {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	log.Fatalf("project root not found from %s", cwd)
	return ""
}

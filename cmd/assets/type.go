package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mechiko/utility"
)

type Assets struct {
	mutex sync.Mutex
	path  string
	jpg   map[string][]byte
	json  map[string][]byte
}

func New(path string) (*Assets, error) {
	if !utility.PathOrFileExists(path) {
		return nil, fmt.Errorf("%s not found", path)
	}
	a := &Assets{
		path: path,
		jpg:  make(map[string][]byte),
		json: make(map[string][]byte),
	}
	err := a.load()
	if err != nil {
		return nil, fmt.Errorf("assets new error %w", err)
	}
	return a, nil
}

func (a *Assets) Jpg(name string) (b []byte, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if name == "" {
		return nil, fmt.Errorf("assets jpg name %s is empty", name)
	}
	name = strings.ToLower(name)
	byteJpg, ok := a.jpg[name]
	if !ok {
		return nil, fmt.Errorf("assets jpg %s not found", name)
	}
	b = make([]byte, len(byteJpg))
	copy(b, byteJpg)
	return
}

func (a *Assets) Json(name string) (b []byte, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if name == "" {
		return nil, fmt.Errorf("assets json name %s is empty", name)
	}
	name = strings.ToLower(name)
	byteJson, ok := a.json[name]
	if !ok {
		return nil, fmt.Errorf("assets json %s not found", name)
	}
	b = make([]byte, len(byteJson))
	copy(b, byteJson)
	return
}

func (a *Assets) load() (err error) {
	entries, err := os.ReadDir(a.path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, file := range entries {
		if !file.IsDir() {
			name := strings.ToLower(file.Name())
			ext := filepath.Ext(name)
			if len(ext) > 3 {
				ext = ext[1:]
			}
			base := name[:len(name)-len(filepath.Ext(name))]
			switch ext {
			case "jpg":
				contentBytes, err := os.ReadFile(filepath.Join(a.path, file.Name()))
				if err != nil {
					return fmt.Errorf("Error reading file: %v", err)
				}
				// Convert the byte slice to a string
				a.jpg[base] = contentBytes
			case "json":
				contentBytes, err := os.ReadFile(filepath.Join(a.path, file.Name()))
				if err != nil {
					return fmt.Errorf("Error reading file: %v", err)
				}
				// Convert the byte slice to a string
				a.json[base] = contentBytes
			}
		}
	}
	return nil
}

package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"pdfimporter/domain"
	"strings"
	"sync"

	"github.com/mechiko/utility"
)

type Assets struct {
	mutex sync.Mutex
	path  string
	jpg   map[string][]byte
	json  map[string][]byte
	// соответствие имени шаблона в описании с именем файла шаблона
	templateNames map[string]string
}

func New(path string) (*Assets, error) {
	if !utility.PathOrFileExists(path) {
		return nil, fmt.Errorf("%s not found", path)
	}
	a := &Assets{
		path:          path,
		jpg:           make(map[string][]byte),
		json:          make(map[string][]byte),
		templateNames: make(map[string]string),
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

func (a *Assets) Templates() (out []string, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	out = make([]string, 0, len(a.templateNames))
	for key := range a.templateNames {
		out = append(out, key)
	}
	return out, nil
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
			if len(ext) > 1 {
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
				out, err := domain.NewMarkTemplate(contentBytes)
				if out == nil || err != nil {
					return fmt.Errorf("new marktemplate error %w", err)
				}
				if out.Name == "" {
					return fmt.Errorf("new marktemplate name empty")
				}
				if _, ok := a.templateNames[out.Name]; ok {
					return fmt.Errorf("marktemplate %s alredy present", out.Name)
				}
				a.templateNames[out.Name] = base
			}
		}
	}
	return nil
}

// находим шаблон по имени в описании
func (a *Assets) Template(name string) (out *domain.MarkTemplate, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if name == "" {
		return nil, fmt.Errorf("assets json name %s is empty", name)
	}
	// name = strings.ToLower(name)
	nameJson, exist := a.templateNames[name]
	if !exist {
		return nil, fmt.Errorf("assets template %s not found", name)
	}
	nameJson = strings.ToLower(nameJson)
	byteJson, ok := a.json[nameJson]
	if !ok {
		return nil, fmt.Errorf("assets json %s not found", nameJson)
	}
	copyByte := make([]byte, len(byteJson))
	copy(copyByte, byteJson)
	out, err = domain.NewMarkTemplate(copyByte)
	return out, err
}

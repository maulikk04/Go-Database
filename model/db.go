package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
	log     Logger
}

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}

	if options != nil {
		opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.DEBUG))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err == nil {
		return &driver, nil
	}

	opts.Logger.Debug("Creating database at '%s'..\n", dir)
	return &driver, os.MkdirAll(dir, 0755)

}

func (d *Driver) Write(collection, id string, v map[string]interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection - no place to save the record")
	}
	if id == "" {
		return fmt.Errorf("missing id -unable to save record")
	}
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, id+".json")
	tmppath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	b = append(b, byte('\n'))
	if err := os.WriteFile(tmppath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmppath, fnlPath)
}

func (d *Driver) Read(collection, id string) (map[string]interface{}, error) {
	if collection == "" {
		return nil, fmt.Errorf("missing collection - unable to read")
	}
	if id == "" {
		return nil, fmt.Errorf("missing id -unable to read record")
	}

	record := filepath.Join(d.dir, collection, id)

	if _, err := stat(record); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return nil, err
	}

	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func (d *Driver) ReadAll(collection string) ([]map[string]interface{}, error) {

	if collection == "" {
		return nil, fmt.Errorf("missing collection - unable to read")
	}

	dir := filepath.Join(d.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := os.ReadDir(dir)

	var records []map[string]interface{}

	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		var v map[string]interface{}

		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}

		records = append(records, v)
	}
	return records, nil
}

func (d *Driver) Delete(collection, id string) error {
	path := filepath.Join(collection, id)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory named %v", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")

	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

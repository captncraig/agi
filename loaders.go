// Package sci contains core types and interfaces for Sierra Creative Interpreter tool
package agi

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Loader interface {
	GetFile(name string) ([]byte, error)
	ListDIRs() ([]string, error)
}

func NewFromDir(base string) Loader {
	return &cachingLoader{
		cache: map[string][]byte{},
		inner: directory(base),
	}
}

type cachingLoader struct {
	cache map[string][]byte
	inner Loader
}

func (c *cachingLoader) GetFile(name string) ([]byte, error) {
	if d, ok := c.cache[name]; ok {
		return d, nil
	}
	return c.inner.GetFile(name)
}

func (c *cachingLoader) ListDIRs() ([]string, error) {
	return c.inner.ListDIRs()
}

type directory string

func (d directory) GetFile(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(string(d), name))
}

func (d directory) ListDIRs() ([]string, error) {
	infos, err := ioutil.ReadDir(string(d))
	dirs := []string{}
	if err != nil {
		return nil, err
	}
	for _, f := range infos {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), "DIR") {
			dirs = append(dirs, strings.TrimSuffix(f.Name(), "DIR"))
		}
	}
	return dirs, nil
}

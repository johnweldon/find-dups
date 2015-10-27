package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var root string

func init() {
	flag.StringVar(&root, "r", ".", "where to start")
}

type meta struct {
	Store map[int64][]string
	Hits  map[int64]interface{}
}

func main() {
	flag.Parse()
	m := meta{Store: make(map[int64][]string), Hits: make(map[int64]interface{})}
	err := filepath.Walk(root, makeWalker(&m))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(-1)
	}

	for k := range m.Hits {
		matches := m.Store[k]
		dups := map[string][]string{}
		for _, v := range matches {
			first, err := firstBytes(v, k)
			if err != nil {
				continue
			}
			hash := fmt.Sprintf("%x", md5.Sum(first))
			paths := dups[hash]
			dups[hash] = append(paths, v)
		}

		for _, v := range dups {
			if len(v) > 1 {
				fmt.Fprintf(os.Stdout, "Potential duplicates:\n")
				for _, p := range v {
					fmt.Fprintf(os.Stdout, "  %s\n", p)
				}
			}
		}
	}
}

func makeWalker(m *meta) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if ignore(path) {
			return nil
		}
		key := info.Size()
		bucket := m.Store[key]
		if len(bucket) > 0 {
			m.Hits[key] = struct{}{}
		}
		m.Store[key] = append(bucket, path)
		return nil
	}
}

func firstBytes(path string, max int64) ([]byte, error) {
	sz := 100
	if int64(sz) > max {
		sz = int(max)
	}
	buf := make([]byte, sz)
	f, err := os.Open(path)
	if err != nil {
		return buf, err
	}
	n, err := f.Read(buf)
	if n != sz {
		return buf, fmt.Errorf("unexpected read length")
	}
	return buf, nil
}

func ignore(path string) bool {
	return strings.Contains(path, ".git")
}

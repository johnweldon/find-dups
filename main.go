package main

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.SetPrefix(" |LOG|> ")
}

type meta struct {
	Store map[int64][]os.FileInfo
	Hits  map[int64]interface{}
}

func main() {
	log.Printf("Begin")
	m := meta{Store: make(map[int64][]os.FileInfo), Hits: make(map[int64]interface{})}
	err := filepath.Walk("..", makeWalker(&m))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("%+v\n", m)

	for k := range m.Hits {
		matches := m.Store[k]
		log.Printf("Size %d: matches: %d", k, len(matches))
		for _, v := range matches {
			log.Printf("     %s", v.Name())
		}
	}
}

func makeWalker(m *meta) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		key := info.Size()
		bucket := m.Store[key]
		if len(bucket) > 0 {
			m.Hits[key] = struct{}{}
		}
		m.Store[key] = append(bucket, info)
		return nil
	}
}

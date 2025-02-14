package repo

import (
	"encoding/gob"
	"log"
	"os"
	path2 "path"
)

type FileRepo[K interface{ string }, V any] struct {
	memRepository MemRepo[K, V]
	path          string
}

// NewFileRepos 创建文件存储
func NewFileRepos[K interface{ string }, V any](path string) *FileRepo[K, V] {
	repo := &FileRepo[K, V]{
		memRepository: MemRepo[K, V]{items: make(map[K]V)},
		path:          path,
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := repo.loadAll()
		if err != nil {
			log.Fatalf("Failed to load data from file %s: %v", path, err)
		}
	} else {
		// 如果文件不存在，检查目录是否存在
		dir := path[:len(path)-len(path2.Base(path))]
		if dir == "" {
			dir = "."
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}
	}
	return repo
}

// writeBack 将数据写回文件
func (f *FileRepo[K, V]) writeBack(_ K, _ V) error {
	log.Print("writeBack to file")
	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		log.Printf("Closing file %s", f.path)
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file %s: %v", f.path, err)
		} else {
			log.Printf("Closed file %s", f.path)
		}
	}(file)

	encoder := gob.NewEncoder(file)
	return encoder.Encode(f.memRepository.items)
}

// loadAll 从文件加载所有数据
func (f *FileRepo[K, V]) loadAll() error {
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file %s: %v", f.path, err)
		}
	}(file)

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&f.memRepository.items)
}

func (f *FileRepo[K, V]) Get(key K) (value V, found bool) {
	return f.memRepository.Get(key)
}

func (f *FileRepo[K, V]) Set(key K, data V) {
	f.memRepository.Set(key, data)
	err := f.writeBack(key, data)
	if err != nil {
		log.Printf("Failed to write back data to file %s: %v", f.path, err)
	}
}

func (f *FileRepo[K, V]) Delete(key K) bool {
	deleted := f.memRepository.Delete(key)
	if deleted {
		err := f.writeBack(key, *new(V)) // Ensure consistency after deletion
		if err != nil {
			log.Printf("Failed to update file after deletion in %s: %v", f.path, err)
		}
	}
	return deleted
}

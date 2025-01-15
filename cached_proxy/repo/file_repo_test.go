package repo

import (
	"os"
	"testing"
)

func TestFileRepo(t *testing.T) {
	filePath := "test_repo.gob"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove file %s: %v", name, err)
		}
	}(filePath) // 清理测试文件

	repo := NewFileRepos[string, string](filePath)

	// 测试 Set 和 Get
	repo.Set("key1", "value1")
	value, found := repo.Get("key1")
	if !found || value != "value1" {
		t.Fatalf("Expected to get 'value1', got '%s'", value)
	}

	// 测试写入文件后重加载
	repo.Set("key2", "value2")
	newRepo := NewFileRepos[string, string](filePath)
	value, found = newRepo.Get("key1")
	if !found || value != "value1" {
		t.Fatalf("Expected to get 'value1' after reload, got '%s'", value)
	}
	value, found = newRepo.Get("key2")
	if !found || value != "value2" {
		t.Fatalf("Expected to get 'value2' after reload, got '%s'", value)
	}

	// 测试删除
	deleted := newRepo.Delete("key1")
	if !deleted {
		t.Fatalf("Expected key1 to be deleted")
	}
	_, found = newRepo.Get("key1")
	if found {
		t.Fatalf("Expected key1 to be absent after deletion")
	}
}

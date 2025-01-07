package repo

import (
	"testing"
)

func TestMemRepo(t *testing.T) {
	// 测试 Set 方法
	t.Run("Set success", func(t *testing.T) {
		repo := NewMemRepo[string, int]()
		repo.items = make(map[string]int) // 初始化内部 map
		repo.Set("key1", 100)
		repo.Set("key2", 200)
		if len(repo.items) != 2 {
			t.Errorf("期望长度为 2，实际长度为 %d", len(repo.items))
		}
	})

	// 测试 Get 方法
	t.Run("Get success", func(t *testing.T) {
		repo := NewMemRepo[string, int]()
		repo.items = make(map[string]int) // 初始化内部 map
		repo.Set("key1", 100)
		repo.Set("key2", 200)
		value, found := repo.Get("key1")
		if !found {
			t.Errorf("期望找到键 'key1'，但未找到")
		}
		if value != 100 {
			t.Errorf("期望值为 100，实际值为 %d", value)
		}
	})

	// 测试 Get 方法（不存在的键）
	t.Run("Get not found", func(t *testing.T) {
		repo := NewMemRepo[string, int]()
		repo.items = make(map[string]int) // 初始化内部 map
		_, found := repo.Get("key3")
		if found {
			t.Errorf("期望未找到键 'key3'，但却找到")
		}
	})

	// 覆盖值测试
	t.Run("Set override", func(t *testing.T) {
		repo := NewMemRepo[string, int]()
		repo.items = make(map[string]int) // 初始化内部 map
		repo.Set("key1", 100)
		repo.Set("key1", 200)
		value, found := repo.Get("key1")
		if !found {
			t.Errorf("期望找到键 'key1'，但未找到")
		}
		if value != 200 {
			t.Errorf("期望值为 200，实际值为 %d", value)
		}
	})

	// 测试 Delete 方法
	t.Run("Delete success", func(t *testing.T) {
		repo := NewMemRepo[string, int]()
		repo.items = make(map[string]int) // 初始化内部 map
		repo.Set("key1", 100)
		repo.Set("key2", 200)
		repo.Delete("key1")
		_, found := repo.Get("key1")
		if found {
			t.Errorf("期望未找到键 'key1'，但却找到")
		}
		if len(repo.items) != 1 {
			t.Errorf("期望长度为 1，实际长度为 %d", len(repo.items))
		}
	})

}

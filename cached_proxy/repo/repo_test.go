package repo

import (
	"testing"
)

func TestMemRepo(t *testing.T) {
	// 创建一个 MemRepo 实例
	repo := NewMemRepo[string, int]()
	repo.items = make(map[string]int) // 初始化内部 map

	// 测试 Set 方法
	repo.Set("key1", 100)
	repo.Set("key2", 200)

	// 测试 Get 方法（存在的键）
	value, found := repo.Get("key1")
	if !found {
		t.Errorf("期望找到键 'key1'，但未找到")
	}
	if value != 100 {
		t.Errorf("期望值为 100，实际值为 %d", value)
	}

	// 测试 Get 方法（另一个存在的键）
	value, found = repo.Get("key2")
	if !found {
		t.Errorf("期望找到键 'key2'，但未找到")
	}
	if value != 200 {
		t.Errorf("期望值为 200，实际值为 %d", value)
	}

	// 测试 Get 方法（不存在的键）
	_, found = repo.Get("key3")
	if found {
		t.Errorf("期望未找到键 'key3'，但却找到")
	}

	// 覆盖值测试
	repo.Set("key1", 300)
	value, found = repo.Get("key1")
	if !found {
		t.Errorf("期望找到键 'key1'，但未找到")
	}
	if value != 300 {
		t.Errorf("期望值为 300，实际值为 %d", value)
	}
}

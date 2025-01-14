package cache

import (
	"cached_proxy/executor"
	"testing"
	"time"
)

func TestPersonalInformationService_GetInfo(t *testing.T) {
	exec := executor.NewWorkerPool(1)
	exec.Run()
	checker := NewIntervalStatusChecker[string](2*time.Second, 3*time.Second)
	onUpdater := func(studentID string) (value *string, update bool) {
		v := "updated info"
		return &v, true
	}
	service := NewPersonalInformationService(exec, checker, onUpdater)

	// Test cache not found
	t.Run("Cache not found", func(t *testing.T) {
		_, err := service.GetInfo("student1")
		if err == nil || err.Error() != "cache not found" {
			t.Errorf("expected cache not found error, got %v", err)
		}
	})

	// Test cache expired
	t.Run("Cache expired", func(t *testing.T) {
		service.(*AbsInfoService[string]).setData("student1", &cacheItem[string]{
			updateAt: time.Now().Add(-4 * time.Second),
			submitAt: time.Now().Add(-4 * time.Second),
			data:     "old info",
		})
		_, err := service.GetInfo("student1")
		if err == nil || err.Error() != "cache expired" {
			t.Errorf("expected cache expired error, got %v", err)
		}
	})

	// Test cache valid
	t.Run("Cache valid", func(t *testing.T) {
		service.(*AbsInfoService[string]).setData("student1", &cacheItem[string]{
			updateAt: time.Now(),
			submitAt: time.Now(),
			data:     "valid info",
		})
		data, err := service.GetInfo("student1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if *data != "valid info" {
			t.Errorf("expected valid info, got %v", *data)
		}
	})

	// Test cache updating
	t.Run("Cache updating", func(t *testing.T) {
		service.(*AbsInfoService[string]).setData("student1", &cacheItem[string]{
			updateAt: time.Now().Add(-4 * time.Second),
			submitAt: time.Now().Add(2 * time.Second),
			data:     "updating info",
		})
		_, err := service.GetInfo("student1")
		if err == nil || err.Error() != "cache updating" {
			t.Errorf("expected cache updating error, got %v", err)
		}
	})
}

func TestPublicInformationService_GetInfo(t *testing.T) {
	exec := executor.NewWorkerPool(1)
	exec.Run()
	checker := NewDailyStatusChecker[string](2 * time.Second)
	onUpdater := func(studentID string) (value *string, update bool) {
		v := "updated info"
		return &v, true
	}
	service := NewPublicInformationService(exec, checker, onUpdater)

	// Test cache not found
	t.Run("Cache not found", func(t *testing.T) {
		_, err := service.GetInfo("student1")
		if err == nil || err.Error() != "cache expired" {
			t.Errorf("expected cache expired error, got %v", err)
		}
	})

	// Test cache expired
	t.Run("Cache expired", func(t *testing.T) {
		service.(*AbsInfoService[string]).setData("student1", &cacheItem[string]{
			updateAt: time.Now().Add(-25 * time.Hour),
			submitAt: time.Now().Add(-25 * time.Hour),
			data:     "old info",
		})
		_, err := service.GetInfo("student1")
		if err == nil || err.Error() != "cache expired" {
			t.Errorf("expected cache expired error, got %v", err)
		}
	})

	// Test cache valid
	t.Run("Cache valid", func(t *testing.T) {
		service.(*AbsInfoService[string]).setData("student1", &cacheItem[string]{
			updateAt: time.Now(),
			submitAt: time.Now(),
			data:     "valid info",
		})
		data, err := service.GetInfo("student1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if *data != "valid info" {
			t.Errorf("expected valid info, got %v", *data)
		}
	})

	// Test cache updating
	t.Run("Cache updating", func(t *testing.T) {
		service.(*AbsInfoService[string]).setData("student1", &cacheItem[string]{
			updateAt: time.Now().Add(-25 * time.Hour),
			submitAt: time.Now().Add(2 * time.Second),
			data:     "updating info",
		})
		_, err := service.GetInfo("student1")
		if err == nil || err.Error() != "cache updating" {
			t.Errorf("expected cache updating error, got %v", err)
		}
	})
}

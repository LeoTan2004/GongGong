package cache

import (
	"testing"
	"time"
)

func TestIntervalStatusChecker_StatusOf(t *testing.T) {
	tests := []struct {
		name     string
		updateAt time.Time
		submitAt time.Time
		expected ItemStatus
	}{
		{
			name:     "Update time expired, submit time not expired",
			updateAt: time.Now().Add(-3 * time.Second),
			submitAt: time.Now().Add(-1 * time.Second),
			expected: Updating,
		},
		{
			name:     "Update time not expired",
			updateAt: time.Now().Add(-1 * time.Second),
			submitAt: time.Now().Add(-4 * time.Second),
			expected: Valid,
		},
		{
			name:     "Update time expired, submit time expired",
			updateAt: time.Now().Add(-3 * time.Second),
			submitAt: time.Now().Add(-4 * time.Second),
			expected: Expired,
		},
	}
	statusChecker := NewIntervalStatusChecker[string](2*time.Second, 3*time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			item := &cacheItem[string]{
				updateAt: tt.updateAt,
				submitAt: tt.submitAt,
			}
			if got := statusChecker.StatusOf(item); got != tt.expected {
				t.Errorf("StatusOf() = %v, expected %v", got, tt.expected)
			}
		})
	}
	t.Run("Item is nil", func(t *testing.T) {
		if got := statusChecker.StatusOf(nil); got != NotFound {
			t.Errorf("StatusOf() = %v, expected false", got)
		}
	})
}

func TestDailyStatusChecker_StatusOf(t *testing.T) {
	tests := []struct {
		name     string
		updateAt time.Time
		submitAt time.Time
		expected ItemStatus
	}{
		{
			name:     "Update time expired, submit time not expired",
			updateAt: time.Now().Add(-25 * time.Hour),
			submitAt: time.Now().Add(-1 * time.Second),
			expected: Updating,
		},
		{
			name:     "Update time not expired",
			updateAt: time.Now().Add(-1 * time.Second),
			submitAt: time.Now().Add(-4 * time.Second),
			expected: Valid,
		},
		{
			name:     "Update time expired, submit time expired",
			updateAt: time.Now().Add(-25 * time.Hour),
			submitAt: time.Now().Add(-4 * time.Second),
			expected: Expired,
		},
	}
	statusChecker := NewDailyStatusChecker[string](2 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			item := &cacheItem[string]{
				updateAt: tt.updateAt,
				submitAt: tt.submitAt,
			}
			if got := statusChecker.StatusOf(item); got != tt.expected {
				t.Errorf("StatusOf() = %v, expected %v", got, tt.expected)
			}
		})
	}
	t.Run("Item is nil", func(t *testing.T) {
		if got := statusChecker.StatusOf(nil); got != NotFound {
			t.Errorf("StatusOf() = %v, expected false", got)
		}
	})
}

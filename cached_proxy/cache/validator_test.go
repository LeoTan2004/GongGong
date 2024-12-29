package cache

import (
	"testing"
	"time"
)

func TestDefaultItemValidator_Valid(t *testing.T) {
	tests := []struct {
		name     string
		updateAt time.Time
		submitAt time.Time
		expected bool
	}{
		{
			name:     "Update time expired, submit time not expired",
			updateAt: time.Now().Add(-3 * time.Second),
			submitAt: time.Now().Add(-1 * time.Second),
			expected: true,
		},
		{
			name:     "Update time not expired",
			updateAt: time.Now().Add(-1 * time.Second),
			submitAt: time.Now().Add(-4 * time.Second),
			expected: true,
		},
	}
	validator := &DefaultItemValidator[string]{
		updateExpireAt: 2 * time.Second,
		submitExpireAt: 3 * time.Second,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			item := &cacheItem[string]{
				updateAt: tt.updateAt,
				submitAt: tt.submitAt,
			}
			if got := validator.Valid(item); got != tt.expected {
				t.Errorf("Valid() = %v, expected %v", got, tt.expected)
			}
		})
	}
	t.Run("Item is nil", func(t *testing.T) {
		if got := validator.Valid(nil); got != false {
			t.Errorf("Valid() = %v, expected false", got)
		}
	})
}

package utils

import (
	"regexp"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// UUID format: xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`
	matched, _ := regexp.MatchString(uuidRegex, uuid)
	if !matched {
		t.Fatalf("Generated UUID %s does not match the expected format", uuid)
	}
}

func TestGenerateMultipleUUIDs(t *testing.T) {
	uuidSet := make(map[string]struct{})
	count := 1000

	for i := 0; i < count; i++ {
		uuid, err := GenerateUUID()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		// Check for uniqueness
		if _, exists := uuidSet[uuid]; exists {
			t.Fatalf("Duplicate UUID found: %s", uuid)
		}
		uuidSet[uuid] = struct{}{}
	}
}

package common

import (
	"testing"
)

type DummyModel struct {
	ID   uint
	Name string
}

func TestOperations_FindByID_RejectsNonNumericID(t *testing.T) {
	ops := NewGormOperations[DummyModel](nil)

	testCases := []string{
		"1 OR 1=1",
		"1; DROP TABLE users;--",
		"abc",
		"",
		"-5",
		"1.5",
	}

	for _, tc := range testCases {
		_, err := ops.FindByID(tc)
		if err == nil {
			t.Errorf("expected error for non-numeric ID %q, got nil", tc)
		}
	}
}

func TestOperations_Delete_RejectsNonNumericID(t *testing.T) {
	ops := NewGormOperations[DummyModel](nil)

	testCases := []string{
		"1 OR 1=1",
		"1; DROP TABLE users;--",
		"abc",
		"",
		"-5",
		"1.5",
	}

	for _, tc := range testCases {
		err := ops.Delete(tc)
		if err == nil {
			t.Errorf("expected error for non-numeric ID %q, got nil", tc)
		}
	}
}

func TestOperations_DeleteMany_EmptyIDs(t *testing.T) {
	ops := NewGormOperations[DummyModel](nil)

	count, err := ops.DeleteMany([]uint{})
	if err != nil {
		t.Fatalf("expected nil error for empty IDs, got: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 deleted rows, got: %d", count)
	}

	count, err = ops.DeleteMany(nil)
	if err != nil {
		t.Fatalf("expected nil error for nil IDs, got: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 deleted rows, got: %d", count)
	}
}

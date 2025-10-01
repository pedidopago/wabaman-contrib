package fbgraph

import "testing"

func TestCompareGraphAPIVersions(t *testing.T) {
	v1 := "v15.0"
	v2 := "v16.0"

	result, err := CompareGraphAPIVersions(v1, v2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != -1 {
		t.Errorf("Expected %s < %s, got %d", v1, v2, result)
	}

	v1 = "v16.0"
	v2 = "v15.0"
	result, err = CompareGraphAPIVersions(v1, v2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected %s > %s, got %d", v1, v2, result)
	}

	v1 = "v15.1"
	v2 = "v15.0"
	result, err = CompareGraphAPIVersions(v1, v2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected %s == %s, got %d", v1, v2, result)
	}

	v1 = "v15.0"
	v2 = "v15.1"
	result, err = CompareGraphAPIVersions(v1, v2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != -1 {
		t.Errorf("Expected %s < %s, got %d", v1, v2, result)
	}

	v1 = "v"
	v2 = "v15.0"
	result, err = CompareGraphAPIVersions(v1, v2)
	if err == nil {
		t.Errorf("Expected error, got %d", result)
	}

	v1 = "vA"
	v2 = "v15.0"
	result, err = CompareGraphAPIVersions(v1, v2)
	if err == nil {
		t.Errorf("Expected error, got %d", result)
	}

	v1 = "v15.0"
	v2 = "v14.999"
	result, err = CompareGraphAPIVersions(v1, v2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected %s < %s, got %d", v1, v2, result)
	}
}

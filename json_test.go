package jchain

import (
	_ "embed"
	"strings"
	"testing"
)

//go:embed testdata/test.json
var embeddedString string

func TestJson(t *testing.T) {
	t.Logf("Successfully embedded string: %s", embeddedString)
	parsed := Parse(embeddedString)
	err := parsed.Error()
	if err != nil {
		t.Fatal(err)
	}
	val, err := parsed.
		Get("menu").
		Get("items").
		Index(3).
		Get("id").
		String()
	if err != nil {
		t.Fatal(err)
	}
	if val != "ZoomIn" {
		t.Fatalf("expected ZoomIn, got %s", val)
	}
}

func TestLargeInteger(t *testing.T) {
	// MaxInt64 is 9223372036854775807
	// We test with MaxInt64 + 1: 9223372036854775808
	jsonStr := `{"big": 9223372036854775808}`
	parsed := Parse(jsonStr)
	if parsed.Error() != nil {
		t.Fatal(parsed.Error())
	}

	val := parsed.Get("big")
	if val.kind != Int {
		t.Fatalf("expected Int kind for large int, got %v", val.kind)
	}

	uVal, err := val.Uint64()
	if err != nil {
		t.Fatal(err)
	}

	expected := uint64(9223372036854775808)
	if uVal != expected {
		t.Errorf("Parsed uint64 value mismatch: got %d, expected %d", uVal, expected)
	}
}

func TestMaxDepth(t *testing.T) {
	// Nested 3 levels: {"a": {"b": {"c": 1}}}
	jsonStr := `{"a": {"b": {"c": 1}}}`

	// Limit 2 -> Should fail
	parsed := ParseWithLimit(jsonStr, 2)
	if parsed.Error() == nil {
		t.Fatal("expected error for exceeding max depth")
	}
	if !strings.Contains(parsed.Error().Error(), "Maximum depth exceeded") {
		t.Fatalf("unexpected error message: %v", parsed.Error())
	}

	// Limit 3 -> Should pass (or 4 depending on how we count)
	// Depth:
	// 1: { (root)
	// 2:   "a": {
	// 3:     "b": {
	// 4:       "c": 1

	// Wait, let's check counting logic.
	// parseObject -> depth check (>= max). depth++.
	// call 1 (root): depth 0 -> 1.
	// call 2 ("a"): depth 1 -> 2.
	// call 3 ("b"): depth 2 -> 3.
	// call 4 ("c"): not an object, but parseValue calls parseNumber... parseNumber doesn't inc depth.
	// If limit is 2.
	// root (depth 0): check 0 >= 2 (false). depth -> 1.
	// "a" val is object. parseValue -> parseObject.
	// "a" (depth 1): check 1 >= 2 (false). depth -> 2.
	// "b" val is object. parseValue -> parseObject.
	// "b" (depth 2): check 2 >= 2 (true) -> ERROR.

	// So limit 2 allows 2 levels of objects?
	// root is level 1?
	// let's verifying behavior.

	parsedPass := ParseWithLimit(jsonStr, 4)
	if parsedPass.Error() != nil {
		t.Fatalf("expected pass for limit 4, got error: %v", parsedPass.Error())
	}
}

func TestErrorMessage(t *testing.T) {
	jsonStr := `
{
	"key": "value",
	"broken": [
		1, 2, 
		error_here
	]
}`
	parsed := Parse(jsonStr)
	if parsed.Error() == nil {
		t.Fatal("expected error for invalid json")
	}
	errMsg := parsed.Error().Error()
	// Line 5, approx column 3 or so. "error_here" starts at tab + error_here?
	// Previous line ends with comma.
	// Line 5: \t\terror_here
	// It should fail at 'e', which is invalid token start.

	t.Logf("Error message: %s", errMsg)
	if !strings.Contains(errMsg, "line 6") {
		t.Errorf("expected error to mention line 6, got: %s", errMsg)
	}
}

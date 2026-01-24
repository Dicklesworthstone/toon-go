package toon

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestAvailable(t *testing.T) {
	// This test assumes tru is available in the test environment
	// If not available, it will skip rather than fail
	if !Available() {
		t.Skip("tru binary not available - skipping availability test")
	}

	path, err := TruPath()
	if err != nil {
		t.Errorf("TruPath() returned error when Available() is true: %v", err)
	}
	if path == "" {
		t.Error("TruPath() returned empty path when Available() is true")
	}
}

func TestEncode_SimpleObject(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	data := map[string]any{"name": "Alice", "age": 30}
	result, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	if !strings.Contains(result, "name: Alice") {
		t.Errorf("Expected output to contain 'name: Alice', got: %s", result)
	}
	if !strings.Contains(result, "age: 30") {
		t.Errorf("Expected output to contain 'age: 30', got: %s", result)
	}
}

func TestEncode_Array(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	data := []int{1, 2, 3, 4, 5}
	result, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Array should contain the values
	if result == "" {
		t.Error("Encode() returned empty result for array")
	}
}

func TestEncode_NestedObject(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	data := map[string]any{
		"user": map[string]any{
			"name":  "Bob",
			"email": "bob@example.com",
		},
		"active": true,
	}
	result, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	if !strings.Contains(result, "Bob") {
		t.Errorf("Expected output to contain 'Bob', got: %s", result)
	}
}

func TestEncode_EmptyObject(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	data := map[string]any{}
	result, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Empty object should produce some valid output
	if result == "" {
		t.Log("Note: Empty object produced empty TOON output (this may be valid)")
	}
}

func TestEncode_SpecialCharacters(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	data := map[string]any{
		"message":  "Hello, \"World\"!",
		"path":     "/home/user/file.txt",
		"unicode":  "日本語 emoji: 🎉",
		"newlines": "line1\nline2\nline3",
		"tabs":     "col1\tcol2\tcol3",
	}
	result, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Verify we get some output
	if result == "" {
		t.Error("Encode() returned empty result for special characters")
	}
}

func TestDecode_SimpleObject(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	toonInput := `name: Alice
age: 30`

	var result map[string]any
	err := Decode(toonInput, &result)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	if result["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got: %v", result["name"])
	}
	if result["age"] != float64(30) { // JSON numbers are float64
		t.Errorf("Expected age=30, got: %v", result["age"])
	}
}

func TestDecode_Array(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	// Simple array format
	toonInput := `[3]: 1, 2, 3`

	var result []any
	err := Decode(toonInput, &result)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 elements, got: %d", len(result))
	}
}

func TestRoundtrip_SimpleObject(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	original := map[string]any{
		"name":   "Test",
		"count":  42,
		"active": true,
	}

	// Encode to TOON
	toonStr, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode back
	var decoded map[string]any
	err = Decode(toonStr, &decoded)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	// Compare
	if decoded["name"] != original["name"] {
		t.Errorf("name mismatch: %v != %v", decoded["name"], original["name"])
	}
	if decoded["count"] != float64(original["count"].(int)) {
		t.Errorf("count mismatch: %v != %v", decoded["count"], original["count"])
	}
	if decoded["active"] != original["active"] {
		t.Errorf("active mismatch: %v != %v", decoded["active"], original["active"])
	}
}

func TestRoundtrip_NestedData(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	original := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice", "age": 30},
			map[string]any{"name": "Bob", "age": 25},
		},
		"metadata": map[string]any{
			"count":   2,
			"version": "1.0",
		},
	}

	// Encode to TOON
	toonStr, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode back
	var decoded map[string]any
	err = Decode(toonStr, &decoded)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	// Verify structure is preserved
	users, ok := decoded["users"].([]any)
	if !ok {
		t.Fatalf("users is not an array: %T", decoded["users"])
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got: %d", len(users))
	}
}

func TestDetectFormat_JSON(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{`{"key": "value"}`, FormatJSON},
		{`[1, 2, 3]`, FormatJSON},
		{`{"nested": {"key": "value"}}`, FormatJSON},
		{`[]`, FormatJSON},
		{`{}`, FormatJSON},
		{`  {"spaces": "ok"}  `, FormatJSON},
		{`"hello: world"`, FormatJSON},
		{`123`, FormatJSON},
		{`true`, FormatJSON},
		{`null`, FormatJSON},
	}

	for _, tt := range tests {
		result := DetectFormat(tt.input)
		if result != tt.expected {
			t.Errorf("DetectFormat(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestDetectFormat_TOON(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"key: value", FormatTOON},
		{"items[3]: a, b, c", FormatTOON},
		{"name: Alice\nage: 30", FormatTOON},
		{"users[2]{name,age}:\n  Alice, 30\n  Bob, 25", FormatTOON},
	}

	for _, tt := range tests {
		result := DetectFormat(tt.input)
		if result != tt.expected {
			t.Errorf("DetectFormat(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestDetectFormat_Unknown(t *testing.T) {
	result := DetectFormat("")
	if result != FormatUnknown {
		t.Errorf("DetectFormat('') = %v, want %v", result, FormatUnknown)
	}
}

func TestEncodeWithOptions_KeyFolding(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	data := map[string]any{
		"user": map[string]any{
			"profile": map[string]any{
				"name":  "Alice",
				"email": "alice@example.com",
			},
		},
	}

	opts := EncodeOptions{
		KeyFolding: "safe",
		Indent:     2,
	}

	result, err := EncodeWithOptions(data, opts)
	if err != nil {
		t.Fatalf("EncodeWithOptions() error: %v", err)
	}

	if result == "" {
		t.Error("EncodeWithOptions() returned empty result")
	}
}

func TestError_TruNotFound(t *testing.T) {
	// Temporarily set env to point to non-existent binary
	oldEnv := os.Getenv("TOON_TRU_BIN")
	os.Setenv("TOON_TRU_BIN", "/nonexistent/path/to/tru")
	defer os.Setenv("TOON_TRU_BIN", oldEnv)

	// Also unset PATH lookup by setting invalid TOON_BIN
	oldBin := os.Getenv("TOON_BIN")
	os.Setenv("TOON_BIN", "/also/nonexistent")
	defer os.Setenv("TOON_BIN", oldBin)

	// If tru is in PATH, this test won't trigger the error
	// In that case, we skip
	_, err := Encode(map[string]any{"test": true})
	if err == nil {
		// tru was found via PATH, skip test
		t.Skip("tru found via PATH - cannot test not-found error")
	}

	toonErr, ok := err.(*ToonError)
	if !ok {
		t.Fatalf("Expected *ToonError, got: %T", err)
	}

	if toonErr.Code != ErrCodeTruNotFound {
		t.Errorf("Expected error code %d, got: %d", ErrCodeTruNotFound, toonErr.Code)
	}
}

func TestError_InvalidJSON(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	// Try to decode invalid TOON that would produce invalid JSON
	invalidToon := "this is not valid toon: {{{"

	var result any
	err := Decode(invalidToon, &result)
	// This may or may not error depending on how tru handles it
	// We just verify we don't panic
	_ = err
}

func TestConvert_JSONToTOON(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	jsonInput := `{"name": "Alice", "age": 30}`

	result, format, err := Convert(jsonInput)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}

	if format != FormatJSON {
		t.Errorf("Expected detected format JSON, got: %v", format)
	}

	// Result should be TOON (not start with {)
	if strings.HasPrefix(strings.TrimSpace(result), "{") {
		t.Errorf("Expected TOON output, but got JSON-like: %s", result)
	}
}

func TestConvert_TOONToJSON(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	toonInput := `name: Alice
age: 30`

	result, format, err := Convert(toonInput)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}

	if format != FormatTOON {
		t.Errorf("Expected detected format TOON, got: %v", format)
	}

	// Result should be valid JSON
	var parsed any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Errorf("Result is not valid JSON: %v", err)
	}
}

func TestDecodeToJSON(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	toonInput := `name: Alice
age: 30`

	result, err := DecodeToJSON(toonInput)
	if err != nil {
		t.Fatalf("DecodeToJSON() error: %v", err)
	}

	// Should be valid JSON
	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	if parsed["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got: %v", parsed["name"])
	}
}

func TestDecodeToValue(t *testing.T) {
	if !Available() {
		t.Skip("tru binary not available")
	}

	toonInput := `name: Alice
age: 30`

	result, err := DecodeToValue(toonInput)
	if err != nil {
		t.Fatalf("DecodeToValue() error: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got: %T", result)
	}

	if m["name"] != "Alice" {
		t.Errorf("Expected name='Alice', got: %v", m["name"])
	}
}

// Benchmark tests
func BenchmarkEncode(b *testing.B) {
	if !Available() {
		b.Skip("tru binary not available")
	}

	data := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice", "age": 30, "email": "alice@example.com"},
			map[string]any{"name": "Bob", "age": 25, "email": "bob@example.com"},
			map[string]any{"name": "Charlie", "age": 35, "email": "charlie@example.com"},
		},
		"metadata": map[string]any{
			"count":   3,
			"version": "1.0",
			"page":    1,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	if !Available() {
		b.Skip("tru binary not available")
	}

	toonInput := `users[3]{name,age,email}:
  Alice, 30, alice@example.com
  Bob, 25, bob@example.com
  Charlie, 35, charlie@example.com
metadata.count: 3
metadata.version: 1.0
metadata.page: 1`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result any
		err := Decode(toonInput, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

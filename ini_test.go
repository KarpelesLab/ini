package ini_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/KarpelesLab/ini"
)

func TestBasicParsing(t *testing.T) {
	f := `; this is a test ini file

var1=value1
var2=value2

[section]
var1=value3
var2=value4`

	ini := ini.New()
	err := ini.Load(bytes.NewReader([]byte(f)))
	if err != nil {
		t.Errorf("failed to parse ini: %s", err)
		return
	}

	if v, ok := ini.Get("root", "var1"); !ok || v != "value1" {
		t.Errorf("failed to get value root/var1, read %#v %#v", v, ok)
	}

	if v, ok := ini.Get("section", "var2"); !ok || v != "value4" {
		t.Errorf("failed to get value section/var2, read %#v %#v", v, ok)
	}
}

func TestCommentStyles(t *testing.T) {
	f := `; semicolon comment
var1=value1
# hash comment
var2=value2`

	ini := ini.New()
	err := ini.Load(bytes.NewReader([]byte(f)))
	if err != nil {
		t.Errorf("failed to parse ini: %s", err)
		return
	}

	if v, ok := ini.Get("root", "var1"); !ok || v != "value1" {
		t.Errorf("failed to get value var1, read %#v", v)
	}
	if v, ok := ini.Get("root", "var2"); !ok || v != "value2" {
		t.Errorf("failed to get value var2, read %#v", v)
	}
}

func TestQuotedValues(t *testing.T) {
	f := `var1="quoted value"
var2='single quoted'
var3="value with = sign"
var4=unquoted`

	ini := ini.New()
	err := ini.Load(bytes.NewReader([]byte(f)))
	if err != nil {
		t.Errorf("failed to parse ini: %s", err)
		return
	}

	if v, ok := ini.Get("root", "var1"); !ok || v != "quoted value" {
		t.Errorf("failed to get quoted value, read %#v", v)
	}
	if v, ok := ini.Get("root", "var2"); !ok || v != "single quoted" {
		t.Errorf("failed to get single quoted value, read %#v", v)
	}
	if v, ok := ini.Get("root", "var3"); !ok || v != "value with = sign" {
		t.Errorf("failed to get value with equals sign, read %#v", v)
	}
	if v, ok := ini.Get("root", "var4"); !ok || v != "unquoted" {
		t.Errorf("failed to get unquoted value, read %#v", v)
	}
}

func TestRoundTripEscaping(t *testing.T) {
	ini1 := ini.New()
	ini1.Set("section", "quotes", `value with "quotes"`)
	ini1.Set("section", "newlines", "line1\nline2\rline3\tindented")
	ini1.Set("section", "special", "special=chars;and#some[more]")
	ini1.Set("section", "backslash", `value with \backslash`)

	buf := &bytes.Buffer{}
	_, err := ini1.WriteTo(buf)
	if err != nil {
		t.Errorf("WriteTo failed: %v", err)
		return
	}

	// Parse the written content back
	ini2 := ini.New()
	_, err = ini2.ReadFrom(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Errorf("ReadFrom failed: %v", err)
		return
	}

	// Verify values were preserved
	if v, ok := ini2.Get("section", "quotes"); !ok || v != `value with "quotes"` {
		t.Errorf("Round-trip quotes mismatch, got %#v", v)
	}

	if v, ok := ini2.Get("section", "newlines"); !ok || v != "line1\nline2\rline3\tindented" {
		t.Errorf("Round-trip newlines mismatch, got %#v", v)
	}

	if v, ok := ini2.Get("section", "special"); !ok || v != "special=chars;and#some[more]" {
		t.Errorf("Round-trip special mismatch, got %#v", v)
	}

	if v, ok := ini2.Get("section", "backslash"); !ok || v != `value with \backslash` {
		t.Errorf("Round-trip backslash mismatch, got %#v", v)
	}
}

func TestSet(t *testing.T) {
	ini := ini.New()
	ini.Set("section", "key", "value")

	if v, ok := ini.Get("section", "key"); !ok || v != "value" {
		t.Errorf("Set failed, got %#v %#v", v, ok)
	}

	// Test overwriting
	ini.Set("section", "key", "new value")
	if v, ok := ini.Get("section", "key"); !ok || v != "new value" {
		t.Errorf("Set overwrite failed, got %#v", v)
	}
}

func TestUnset(t *testing.T) {
	ini := ini.New()
	ini.Set("section", "key1", "value1")
	ini.Set("section", "key2", "value2")

	// Test removing a key
	ini.Unset("section", "key1")
	if _, ok := ini.Get("section", "key1"); ok {
		t.Errorf("Unset failed, key still exists")
	}

	// Make sure other keys in the same section still exist
	if v, ok := ini.Get("section", "key2"); !ok || v != "value2" {
		t.Errorf("Unset affected other keys, got %#v %#v", v, ok)
	}

	// Test that removing the last key in a section removes the section
	ini.Unset("section", "key2")
	if ini.HasSection("section") {
		t.Errorf("Removing last key didn't remove section")
	}
}

func TestGetDefault(t *testing.T) {
	ini := ini.New()
	ini.Set("section", "key", "value")

	// Test existing key
	if v := ini.GetDefault("section", "key", "default"); v != "value" {
		t.Errorf("GetDefault for existing key returned %#v, wanted %#v", v, "value")
	}

	// Test non-existing key
	if v := ini.GetDefault("section", "nonexistent", "default"); v != "default" {
		t.Errorf("GetDefault for non-existing key returned %#v, wanted %#v", v, "default")
	}
}

func TestSections(t *testing.T) {
	ini := ini.New()
	ini.Set("section1", "key", "value")
	ini.Set("section2", "key", "value")

	sections := ini.Sections()
	if len(sections) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(sections))
	}

	// Check section names (order not guaranteed)
	found1, found2 := false, false
	for _, s := range sections {
		if s == "section1" {
			found1 = true
		} else if s == "section2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("Sections() didn't return all section names")
	}
}

func TestKeys(t *testing.T) {
	ini := ini.New()
	ini.Set("section", "key1", "value1")
	ini.Set("section", "key2", "value2")

	keys := ini.Keys("section")
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	// Check key names (order not guaranteed)
	found1, found2 := false, false
	for _, k := range keys {
		if k == "key1" {
			found1 = true
		} else if k == "key2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("Keys() didn't return all key names")
	}
}

func TestWriteAndRead(t *testing.T) {
	ini1 := ini.New()
	ini1.Set("root", "key1", "value1")
	ini1.Set("section", "key2", "value2")
	ini1.Set("section", "key with spaces", "value with spaces")

	buf := &bytes.Buffer{}
	if err := ini1.Write(buf); err != nil {
		t.Errorf("Write failed: %v", err)
		return
	}

	// Parse the written content back
	ini2 := ini.New()
	if err := ini2.Load(bytes.NewReader(buf.Bytes())); err != nil {
		t.Errorf("Failed to reload written data: %v", err)
		return
	}

	// Verify values were preserved
	if v, ok := ini2.Get("root", "key1"); !ok || v != "value1" {
		t.Errorf("Round-trip value mismatch for root/key1, got %#v", v)
	}

	if v, ok := ini2.Get("section", "key2"); !ok || v != "value2" {
		t.Errorf("Round-trip value mismatch for section/key2, got %#v", v)
	}

	if v, ok := ini2.Get("section", "key with spaces"); !ok || v != "value with spaces" {
		t.Errorf("Round-trip value mismatch for key with spaces, got %#v", v)
	}
}

func TestReaderFromWriterTo(t *testing.T) {
	ini1 := ini.New()
	ini1.Set("root", "key1", "value1")
	ini1.Set("section", "key2", "value2")
	ini1.Set("section", "key with spaces", "value with spaces")

	buf := &bytes.Buffer{}
	bytesWritten, err := ini1.WriteTo(buf)
	if err != nil {
		t.Errorf("WriteTo failed: %v", err)
		return
	}

	if bytesWritten <= 0 {
		t.Errorf("Expected bytesWritten > 0, got %d", bytesWritten)
	}

	// Parse the written content back
	ini2 := ini.New()
	bytesRead, err := ini2.ReadFrom(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Errorf("ReadFrom failed: %v", err)
		return
	}

	if bytesRead <= 0 {
		t.Errorf("Expected bytesRead > 0, got %d", bytesRead)
	}

	// Verify values were preserved
	if v, ok := ini2.Get("root", "key1"); !ok || v != "value1" {
		t.Errorf("Round-trip value mismatch for root/key1, got %#v", v)
	}

	if v, ok := ini2.Get("section", "key2"); !ok || v != "value2" {
		t.Errorf("Round-trip value mismatch for section/key2, got %#v", v)
	}

	if v, ok := ini2.Get("section", "key with spaces"); !ok || v != "value with spaces" {
		t.Errorf("Round-trip value mismatch for key with spaces, got %#v", v)
	}
}

func TestThreadSafeReadWrite(t *testing.T) {
	ini1 := ini.NewThreadSafe()
	ini1.Set("root", "key1", "value1")
	ini1.Set("section", "key2", "value2")

	buf := &bytes.Buffer{}
	bytesWritten, err := ini1.WriteTo(buf)
	if err != nil {
		t.Errorf("WriteTo failed: %v", err)
		return
	}

	if bytesWritten <= 0 {
		t.Errorf("Expected bytesWritten > 0, got %d", bytesWritten)
	}

	// Parse the written content back
	ini2 := ini.NewThreadSafe()
	bytesRead, err := ini2.ReadFrom(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Errorf("ReadFrom failed: %v", err)
		return
	}

	if bytesRead <= 0 {
		t.Errorf("Expected bytesRead > 0, got %d", bytesRead)
	}

	// Verify values were preserved
	if v, ok := ini2.Get("root", "key1"); !ok || v != "value1" {
		t.Errorf("Round-trip value mismatch for root/key1, got %#v", v)
	}

	if v, ok := ini2.Get("section", "key2"); !ok || v != "value2" {
		t.Errorf("Round-trip value mismatch for section/key2, got %#v", v)
	}
}

func TestErrorCases(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		errorMsg string
	}{
		{
			name:     "Empty section name",
			content:  "[]\nkey=value",
			errorMsg: "empty section name",
		},
		{
			name:     "Empty key name",
			content:  "=value",
			errorMsg: "empty key name",
		},
		{
			name:     "Invalid line",
			content:  "invalid line without equals",
			errorMsg: "invalid format, missing '='",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ini := ini.New()
			err := ini.Load(strings.NewReader(tc.content))

			if err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
				return
			}

			if !strings.Contains(err.Error(), tc.errorMsg) {
				t.Errorf("Expected error containing %q, got %q", tc.errorMsg, err.Error())
			}
		})
	}
}

func TestThreadSafeBasic(t *testing.T) {
	ini := ini.NewThreadSafe()
	ini.Set("section", "key", "initial")

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)

		// Reader goroutine
		go func() {
			defer wg.Done()
			// Read operations
			ini.Get("section", "key")
			ini.HasSection("section")
			ini.Sections()
			ini.Keys("section")
		}()

		// Writer goroutine
		go func() {
			defer wg.Done()
			// Write operations
			ini.Set("concurrent", "key", "value")
			ini.Get("concurrent", "key")
			ini.Unset("concurrent", "key")
		}()
	}

	wg.Wait()
	// If we got here without race detector errors, we're good
}

package ifacecommand

import (
	"testing"
)

// TestNew tests interface command creation
func TestNew(t *testing.T) {
	cmd := New()
	if cmd == nil {
		t.Fatal("New() returned nil")
	}

	if cmd.Name() != "interface" {
		t.Errorf("Expected command name 'interface', got %s", cmd.Name())
	}

	if cmd.Description() != "generate interface from struct" {
		t.Errorf("Expected description 'generate interface from struct', got %s", cmd.Description())
	}
}

// TestParse tests command flag parsing and validation
func TestParse(t *testing.T) {
	cmd := New()

	// Test with required flags
	args := []string{"-pkg", ".", "-type", "TestStruct", "-out", "test.go"}
	err := cmd.Parse(args)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// Test missing type flag (should fail)
	cmd2 := New()
	args2 := []string{"-pkg", "."}
	err2 := cmd2.Parse(args2)
	if err2 == nil {
		t.Error("Parse() should fail when -type flag is missing")
	}
}

// TestExecute tests command execution with error handling
func TestExecute(t *testing.T) {
	cmd := New()

	// This should fail because no valid target is specified
	exitCode := cmd.Execute()
	if exitCode == 0 {
		t.Error("Execute() should return non-zero exit code when no valid target is specified")
	}
}

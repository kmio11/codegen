package sample

import (
	"os"
	"path/filepath"
	"testing"

	ifacecommand "github.com/kmio11/codegen/cmd/interface"
)

// RED TEST: Integration test for end-to-end interface generation
func TestIntegrationInterfaceGeneration(t *testing.T) {
	// Create temporary output file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "user_interface.go")
	
	// Create command
	cmd := ifacecommand.New()
	
	// Test with UserService struct
	args := []string{
		"-pkg", ".",
		"-type", "UserService",
		"-out", outputFile,
		"-name", "UserRepository",
	}
	
	err := cmd.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse command args: %v", err)
	}
	
	// Execute command
	exitCode := cmd.Execute()
	if exitCode != 0 {
		t.Fatalf("Command execution failed with exit code: %d", exitCode)
	}
	
	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}
	
	// Read and verify output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	expectedContent := []string{
		"package sample",
		"type UserRepository interface {",
		"GetUser(id string) (string, error)",
		"CreateUser(id string, name string) error",
		"UpdateUser(id string, name string) error",
		"DeleteUser(id string) error",
		"ListUsers() map[string]string",
	}
	
	contentStr := string(content)
	for _, expected := range expectedContent {
		if !contains(contentStr, expected) {
			t.Errorf("Expected content not found: %s", expected)
			t.Logf("Generated content:\n%s", contentStr)
		}
	}
	
	// Verify private method is not included
	if contains(contentStr, "privateMethod") {
		t.Error("Private method should not be included in interface")
	}
}

// RED TEST: Test interface generation with default name
func TestIntegrationDefaultInterfaceName(t *testing.T) {
	// Create temporary output file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "default_interface.go")
	
	// Create command
	cmd := ifacecommand.New()
	
	// Test with UserService struct (no custom name)
	args := []string{
		"-pkg", ".",
		"-type", "UserService",
		"-out", outputFile,
	}
	
	err := cmd.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse command args: %v", err)
	}
	
	// Execute command
	exitCode := cmd.Execute()
	if exitCode != 0 {
		t.Fatalf("Command execution failed with exit code: %d", exitCode)
	}
	
	// Read and verify output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	contentStr := string(content)
	
	// Should use default interface name
	if !contains(contentStr, "type UserServiceInterface interface {") {
		t.Error("Expected default interface name 'UserServiceInterface' not found")
		t.Logf("Generated content:\n%s", contentStr)
	}
}

// RED TEST: Test error handling for non-existent struct
func TestIntegrationNonExistentStruct(t *testing.T) {
	// Create command
	cmd := ifacecommand.New()
	
	// Test with non-existent struct
	args := []string{
		"-pkg", ".",
		"-type", "NonExistentStruct",
	}
	
	err := cmd.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse command args: %v", err)
	}
	
	// Execute command - should fail
	exitCode := cmd.Execute()
	if exitCode == 0 {
		t.Error("Command should fail for non-existent struct")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   findSubstring(s, substr) != -1
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
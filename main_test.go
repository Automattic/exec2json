package main

import (
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedStatus int
		expectedStdout string
		expectedStderr string
		expectError    bool
	}{
		{
			name:           "Echo command",
			args:           []string{"echo", "Hello, World!"},
			expectedStatus: 0,
			expectedStdout: "Hello, World!\n",
			expectedStderr: "",
			expectError:    false,
		},
		{
			name:           "Non-existent command",
			args:           []string{"non_existent_command"},
			expectedStatus: -1, // The status will not be set for a non-existent command
			expectedStdout: "",
			expectedStderr: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := executeCommand(tc.args)

			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != nil {
				if result.Status != tc.expectedStatus {
					t.Errorf("Expected status %d, got %d", tc.expectedStatus, result.Status)
				}

				if result.Stdout != tc.expectedStdout {
					t.Errorf("Expected stdout %q, got %q", tc.expectedStdout, result.Stdout)
				}

				if result.Stderr != tc.expectedStderr {
					t.Errorf("Expected stderr %q, got %q", tc.expectedStderr, result.Stderr)
				}

				if result.Took <= 0 {
					t.Errorf("Expected 'took' to be greater than 0, got %f", result.Took)
				}

				if len(result.Command) != len(tc.args) {
					t.Errorf("Expected command %v, got %v", tc.args, result.Command)
				}
			} else if !tc.expectError {
				t.Errorf("Expected a result, but got nil")
			}
		})
	}
}

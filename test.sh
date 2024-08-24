#!/bin/bash

set -e

# Check if exec2json executable exists
if [ ! -f "./exec2json" ]; then
    echo "Error: exec2json executable not found. Make sure to build it before running tests."
    exit 1
fi

# Function to run a test case
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_stdout="$3"
    local expected_stderr="$4"
    local expected_status="$5"

    echo "Running test: $test_name"
    output=$(./exec2json $command)
    
    actual_stdout=$(echo "$output" | jq -r '.stdout' | tr -d '\n')
    actual_stderr=$(echo "$output" | jq -r '.stderr' | tr -d '\n')
    actual_status=$(echo "$output" | jq -r '.status')
    actual_command=$(echo "$output" | jq -r '.command | join(" ")')

    # Remove surrounding quotes if present
    actual_stdout="${actual_stdout#\'}"
    actual_stdout="${actual_stdout%\'}"
    actual_stderr="${actual_stderr#\'}"
    actual_stderr="${actual_stderr%\'}"

    if [ "$actual_stdout" != "$expected_stdout" ] || 
       [ "$actual_stderr" != "$expected_stderr" ] || 
       [ "$actual_status" != "$expected_status" ] ||
       [ "$actual_command" != "$command" ]; then
        echo "Test case '$test_name' failed"
        echo "Command: $command"
        echo "Expected stdout: $expected_stdout"
        echo "Actual stdout: $actual_stdout"
        echo "Expected stderr: $expected_stderr"
        echo "Actual stderr: $actual_stderr"
        echo "Expected status: $expected_status"
        echo "Actual status: $actual_status"
        echo "Actual command: $actual_command"
        exit 1
    fi

    echo "Test case '$test_name' passed"
    echo
}

# Test cases
run_test "Echo command" "echo 'Hello, World!'" "Hello, World!" "" "0"
run_test "Echo multiple words" "echo word1 word2 word3" "word1 word2 word3" "" "0"
run_test "Echo with quotes" "echo 'quoted string'" "quoted string" "" "0"
run_test "Command with arguments" "printf '%s %s' arg1 arg2" "arg1 arg2" "" "0"
run_test "Exit with non-zero status" "bash -c 'exit 42'" "" "" "42"
run_test "Command with stderr output" "bash -c 'echo error >&2; echo output'" "output" "error" "0"
run_test "Non-existent command" "nonexistentcommand" "" "nonexistentcommand: command not found" "127"

echo "All tests passed!"

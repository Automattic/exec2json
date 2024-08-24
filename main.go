/**
 * Copyright 2024 Automattic, Inc.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// CommandResult represents the result of executing a command
type CommandResult struct {
	Command []string  `json:"command"`
	Stdout  string    `json:"stdout"`
	Stderr  string    `json:"stderr"`
	Status  int       `json:"status"`
	Took    float64   `json:"took"`
}

func executeCommand(args []string) (*CommandResult, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("you must supply a command to execute")
	}
	command := exec.Command(args[0], args[1:]...)
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	startTime := time.Now()
	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()
	err = command.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}
	outString, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdout: %w", err)
	}
	errString, err := io.ReadAll(stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to read stderr: %w", err)
	}
	err = command.Wait()
	took := time.Now().Sub(startTime).Seconds()
	result := &CommandResult{
		Command: args,
		Stdout:  string(outString),
		Stderr:  string(errString),
		Status:  command.ProcessState.ExitCode(),
		Took:    took,
	}
	return result, err
}

func main() {
	result, err := executeCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	json.NewEncoder(os.Stdout).Encode(result)
	os.Exit(result.Status)
}

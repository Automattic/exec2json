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
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must supply a command to execute")
	}
	command := exec.Command(os.Args[1], os.Args[2:]...)

	stdin, err := command.StdinPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}

	startTime := time.Now()

	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()

	err = command.Start()
	if err != nil {
		panic(err)
	}

	outString, err := io.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	errString, err := io.ReadAll(stderr)
	if err != nil {
		panic(err)
	}

	if err := command.Wait(); err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			panic(err)
		}
	}
	took := time.Now().Sub(startTime).Seconds()

	json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"command": os.Args[1:],
		"stdout":  string(outString),
		"stderr":  string(errString),
		"status":  command.ProcessState.ExitCode(),
		"took":    took,
	})

	os.Exit(command.ProcessState.ExitCode())
}

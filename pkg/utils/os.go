// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SearchExecutable(name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("error, empty name")
	}

	if path, err := exec.LookPath(name); err != nil {
		return "", err
	} else {
		return path, nil
	}

}

func RunFrondendCmd(ctx context.Context, cmdName string, env []string, stdin_msg string) (stdoutMsg, stderrMsg string, exitedCode int, e error) {

	var outMsg bytes.Buffer
	var outErr bytes.Buffer
	var cmd *exec.Cmd

	if len(cmdName) == 0 {
		e = errors.New("error, empty cmd")
		return
	}

	rootCmd := "bash"
	if path, _ := SearchExecutable(rootCmd); len(path) != 0 {
		cmd = exec.CommandContext(ctx, rootCmd, "-c", cmdName)
		goto EXE
	}

	rootCmd = "sh"
	if path, _ := SearchExecutable(rootCmd); len(path) != 0 {
		cmd = exec.CommandContext(ctx, rootCmd, "-c", cmdName)
		goto EXE
	}

	e = errors.New("error, no sh or bash installed")
	return

EXE:

	cmd.Env = append(os.Environ(), env...)

	if len(stdin_msg) != 0 {
		cmd.Stdin = strings.NewReader(stdin_msg)
	}

	cmd.Stdout = &outMsg
	cmd.Stderr = &outErr

	e = cmd.Run()
	if a := strings.TrimSpace(outMsg.String()); len(a) > 0 {
		stdoutMsg = a
	}
	if b := strings.TrimSpace(outErr.String()); len(b) > 0 {
		stderrMsg = b
	}
	exitedCode = cmd.ProcessState.ExitCode()

	return
}

func GetFileList(dirName string) ([]string, error) {
	filelist, e := os.ReadDir(dirName)
	if e != nil {
		return nil, fmt.Errorf("failed to read directory %s, error=%v", dirName, e)
	}

	nameList := []string{}
	for _, item := range filelist {
		if item.IsDir() {
			continue
		}
		nameList = append(nameList, item.Name())
	}
	return nameList, nil
}

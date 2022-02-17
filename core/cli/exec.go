package cli

import (
	"os/exec"
	"runtime"
)

/**
 * Execute command
 */
func ExecOrPanic(cmd string) string {
	output, err := exec.Command(cmd).Output()
	if err != nil {
		panic(err)
	}
	return string(output)
}

/**
 * Execute command inside shell on linux or windows
 */
func ExecShellScriptOrPanic(args string) string {
	var (
		shell string
		cmd   *exec.Cmd
	)

	switch runtime.GOOS {
	case "windows":
		shell = "cmd.exe"
		cmdArgs := []string{"/C", args}
		cmd = exec.Command(shell, cmdArgs...)
		break
	case "linux":
		shell = "bash"
		cmdArgs := []string{"-c", args}
		cmd = exec.Command(shell, cmdArgs...)
		break
	default:
		panic("Unsupported target OS, unable to exec command !")
		break
	}

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(output)
}

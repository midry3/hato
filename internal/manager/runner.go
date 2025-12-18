package manager

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	hasSh bool
)

func RunCmd(command string) error {
	var (
		cmd *exec.Cmd
		env []string
		c   strings.Builder
		f   bool
	)
	env = os.Environ()
	for _, a := range strings.Fields(command) {
		if f {
			c.WriteString(" " + a)
			continue
		}
		if strings.Count(a, "=") == 0 {
			f = true
			c.WriteString(" " + a)
		} else {
			env = append(env, a)
		}
	}
	if hasSh {
		cmd = exec.Command("sh", "-c", c.String())
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", c.String())
	}
	cmd.Env = env
	cmd.Dir, _ = os.Getwd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	_, err := exec.LookPath("sh")
	hasSh = err == nil
	if !hasSh && runtime.GOOS != "windows" {
		fmt.Fprintln(os.Stderr, "\033[31mNOTE\033[0m: Could not find shell, so actions will not be run.")
	}
}

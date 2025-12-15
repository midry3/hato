package manager

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RunCmd(command string) error {
	var (
		cmd *exec.Cmd
		env []string
		c   []string
		f   bool
	)
	env = os.Environ()
	for _, a := range strings.Fields(command) {
		if f {
			c = append(c, a)
			continue
		}
		k := strings.SplitN(a, "=", 2)
		if len(k) == 1 {
			f = true
			c = append(c, a)
		} else {
			env = append(env, a)
		}
	}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", strings.Join(c, " "))
	} else {
		cmd = exec.Command("sh", "-c", strings.Join(c, " "))
	}
	cmd.Env = env
	cmd.Dir, _ = os.Getwd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Env = replaceEnvs(env, os.Environ())

	if err := c.Run(); err != nil {
		return c.ProcessState.ExitCode()
	}
	return c.ProcessState.ExitCode()
}

func replaceEnvs(newEnv Environment, osEnv []string) []string {
	replaced := make([]string, 0, len(osEnv)+len(newEnv))
	for _, v := range osEnv {
		key := strings.SplitN(v, "=", 2)[0]
		newVal, exists := newEnv[key]

		delete(newEnv, key)

		if exists && newVal.NeedRemove {
			continue
		}
		if exists {
			v = key + "=" + newVal.Value
		}
		replaced = append(replaced, v)
	}
	for key, v := range newEnv {
		replaced = append(replaced, key+"="+v.Value)
	}
	return replaced
}

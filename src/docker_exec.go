package dockerImage

import (
	"fmt"
	osExec "os/exec"
)

type dockerExec string

func (exec dockerExec) validateExecutable() error {
	dockerVersionCommand := osExec.Command(string(exec), "-v")
	dockerVersionCommandError := dockerVersionCommand.Run()

	if dockerVersionCommandError == osExec.ErrNotFound {
		return fmt.Errorf("docker executable '%s' not found: %s", exec, dockerVersionCommandError)
	}

	if dockerVersionCommandError != nil {
		return fmt.Errorf("docker version command failed: %s", dockerVersionCommandError)
	}

	return nil
}

func (exec dockerExec) buildContainer(pathToDockerfile string, registry string) error {
	cmd := osExec.Command(string(exec), "build", "-t", registry, pathToDockerfile)

	tagCmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("'build -t' command failed: %s\nOutput:\n%s", err, tagCmdOutput)
	}

	return nil
}

func (exec dockerExec) pushContainer(imageID string, tag string, registry string) error {
	cmd := osExec.Command(string(exec), "tag", imageID, registry)
	tagCmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("'tag' command failed: %s\nOutput:\n%s", err, tagCmdOutput)
	}

	cmd = osExec.Command(string(exec), "push", registry)
	pushCmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("'push' command failed: %s\nOutput:\n%s", err, pushCmdOutput)
	}

	return nil
}

package dockerImage

import (
	"fmt"
	osExec "os/exec"
	"strings"
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

func (exec dockerExec) buildContainer(pathToDockerfile string, tag string) (string, error) {
	cmd := osExec.Command(string(exec), "build", "-t", tag, pathToDockerfile)

	tagCmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("'build -t' command failed: %s\nOutput:\n%s", err, tagCmdOutput)
	}

	getHashCmd := osExec.Command(string(exec), "images", "-q", tag)
	getHashCmdOutput, err := getHashCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("'images -q' command failed: %s\nOutput:\n%s", err, getHashCmdOutput)
	}

	// TODO: is the [:] necessary? the linter doesn't complain
	return strings.TrimSpace(string(getHashCmdOutput[:])), nil
}

func (exec dockerExec) deleteContainer(imageID string) error {
	cmd := osExec.Command(string(exec), "rmi", imageID)
	rmiCmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("'rmi' command failed: %s\nOutput:\n%s", err, rmiCmdOutput)
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

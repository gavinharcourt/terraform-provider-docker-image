package dockerImage

import osExec "os/exec"

type dockerExec string

// TODO: need to validate pathToDockerfile in schema
func (exec dockerExec) buildContainer(pathToDockerfile string, tag string) (string, error) {
	cmd := osExec.Command(string(exec), "build", "-t", tag, pathToDockerfile)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	getHashCmd := osExec.Command(string(exec), "images", "-q", tag)
	getHashCmdOutput, err := getHashCmd.Output()
	if err != nil {
		return "", err
	}

	// TODO: is the [:] necessary? the linter doesn't complain
	return string(getHashCmdOutput[:]), nil
}

func (exec dockerExec) deleteContainer(imageID string) error {
	cmd := osExec.Command(string(exec), "rmi", imageID)
	return cmd.Run()
}

func (exec dockerExec) pushContainer(imageID string, tag string, repository string) error {
	return nil
}

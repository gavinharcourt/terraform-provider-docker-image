package dockerImage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var basePath string
var registryId string
var testImageTag string

func TestMain(m *testing.M) {
	testImageTag = "terraform-provider-docker-image-test"

	checkTestCanRun()
	changeDockerContents()

	os.Exit(m.Run())
}

func checkTestCanRun() {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		panic("AWS_ACCESS_KEY_ID must be set for the tests to run")
	}

	accessSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessSecret == "" {
		panic("AWS_SECRET_ACCESS_KEY must be set for the tests to run")
	}

	registryId = os.Getenv("ECR_REPOSITORY")
	if registryId == "" {
		panic("ECR_REPOSITORY must be set for the tests to run")
	}

	var err error
	basePath, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get base path: %s", err))
	}
	basePath += "/.."
}

func changeDockerContents() {
	packagedDataFile := basePath + "/test/resources/app/packaged_data"

	packagedDataContents := []byte(time.Now().String())
	err := ioutil.WriteFile(packagedDataFile, packagedDataContents, 0644)
	if err != nil {
		panic(fmt.Sprintf("couldn't create packaged_data file: %s", err))
	}
}

func TestTerraformPlan(t *testing.T) {
	planCmd := exec.Command("terraform", "plan")
	planCmd.Dir = basePath + "/test/resources/terraform"

	outBytes, err := planCmd.CombinedOutput()
	out := string(outBytes)
	if err != nil {
		t.Errorf("Failed to run 'terraform plan': %s\nOutput: %s", err, out)
		return
	}

	if !strings.Contains(out, "+ dockerimage_local.test_image") {
		t.Error("'terrform plan' doesn't want to build dockerimage_local resource.")
	}

	if !strings.Contains(out, "+ dockerimage_remote.test_image") {
		t.Error("'terraform plan' doesn't want to build dockerimage_remote resource.")
	}

	if !strings.Contains(out, "2 to add, 0 to change, 0 to destroy.") {
		t.Error("Expected 'terraform plan' to want to add two resources.")
	}

	t.Logf("Actual output: %s", out)
}

func TestTerraformApply(t *testing.T) {
	existingRemoteSha, _ := getRemoteDockerImage()

	applyCmd := exec.Command("terraform", "apply")
	applyCmd.Dir = basePath + "/test/resources/terraform"

	outBytes, err := applyCmd.CombinedOutput()
	out := string(outBytes)
	if err != nil {
		t.Errorf("Failed to run 'terraform apply': %s\nOutput: %s", err, out)
		return
	}

	localDockerImage, err := getLocalDockerImage()
	if err != nil {
		t.Errorf("could not get local docker image ID: %s", err)
		return
	}

	if localDockerImage == "" {
		t.Errorf("'terraform apply' failed to build local docker image")
		return
	}

	remoteDockerImage, err := getRemoteDockerImage()
	if err != nil {
		t.Errorf("could not get remote docker image ID: %s", err)
		return
	}

	if remoteDockerImage == existingRemoteSha {
		t.Errorf("'terraform apply' did not push built docker image; remote has not changed: %s (new) == %s (old)",
			remoteDockerImage, existingRemoteSha)
		return
	}
}

func getLocalDockerImage() (string, error) {
	getHashCmd := exec.Command("docker", "images", "--no-trunc", "-q", testImageTag)
	getHashCmdOutput, err := getHashCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to get image ID: %s\nOutput: %s", err, getHashCmdOutput)
	}

	return strings.TrimSpace(string(getHashCmdOutput)), nil
}

func getRemoteDockerImage() (string, error) {
	getImageInfoCmd := exec.Command("aws", "ecr", "batch-get-image",
		"--registry-id="+registryId, "--repository-name="+testImageTag, "--image-ids=imageTag=latest")
	getImageInfoCmdOutputBytes, err := getImageInfoCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to run batch-get-image: %s\nOutput: %s", err, getImageInfoCmdOutputBytes)
	}

	var getImageInfoCmdOutput interface{}
	json.Unmarshal(getImageInfoCmdOutputBytes, &getImageInfoCmdOutput)

	imageInfo := getImageInfoCmdOutput.(map[string]interface{})["images"].([]interface{})[0]
	imageDigest := imageInfo.(map[string]interface{})["imageId"].(map[string]interface{})["imageDigest"].(string)

	return imageDigest, nil
}

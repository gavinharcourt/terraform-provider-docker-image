package dockerImage

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os/exec"
)

// Provider for local/remote docker
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// docker config
			"docker_executable": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "docker",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"local_docker_image": dataSourceLocalDockerImage(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	dockerExecutable := d.Get("docker_executable").(string)

	// TODO: move this to docker_exec.go
	// check if executable exists
	dockerVersionCommand := exec.Command(dockerExecutable, "-v")
	dockerVersionCommandError := dockerVersionCommand.Run()
	if dockerVersionCommandError == exec.ErrNotFound {
		return nil, fmt.Errorf("docker executable '%s' not found: %s", dockerExecutable, dockerVersionCommandError)
	} else if dockerVersionCommandError != nil {
		return nil, fmt.Errorf("docker version command failed: %s", dockerVersionCommandError)
	}

	client := &Config{
		DockerExecutable: dockerExecutable,
	}

	return client, nil
}

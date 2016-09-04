package dockerImage

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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

	err := dockerExec(dockerExecutable).validateExecutable()
	if err != nil {
		return nil, err
	}

	client := &Config{
		DockerExecutable: dockerExecutable,
	}

	return client, nil
}

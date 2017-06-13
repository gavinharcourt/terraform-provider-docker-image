package dockerImage

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceLocalDockerImage() *schema.Resource {
	return &schema.Resource{
		Create: dataSourceLocalDockerImageCreate,
		Read:   dataSourceLocalDockerImageRead,
		Exists: dataSourceLocalDockerImageExists,
		Delete: dataSourceLocalDockerImageDelete,

		Schema: map[string]*schema.Schema{
			"dockerfile_path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Absolute path to the Dockerfile to build.",
				ForceNew:    true,
			},

			"registry": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remote registry to push the image to.",
				Default:     "",
				Sensitive:   true,
			},
		},
	}
}

func dataSourceLocalDockerImageCreate(d *schema.ResourceData, meta interface{}) error {
	pathToDockerfile := d.Get("dockerfile_path").(string)
	registry := d.Get("registry").(string)

	err := dockerExec(meta.(*Config).DockerExecutable).buildContainer(pathToDockerfile, registry)
	if err != nil {
		return fmt.Errorf("Failed to create local docker image: %s", err)
	}

	d.SetId(registry)
	return nil
}

// following template data provider convention of doing real work in Exists.
// not sure how appropriate it is for a docker image.
func dataSourceLocalDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func dataSourceLocalDockerImageExists(d *schema.ResourceData, meta interface{}) error {
	pathToDockerfile := d.Get("dockerfile_path").(string)
	registry := d.Get("registry").(string)

	err := dockerExec(meta.(*Config).DockerExecutable).buildContainer(pathToDockerfile, registry)
	if err != nil {
		return false, fmt.Errorf("Failed to build local docker image: %s", err)
	}

	d.SetId(registry)
	return nil
}

func dataSourceLocalDockerImageDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[WARN] deleting local docker images is not currently supported, they must be deleted manually.")
	return nil
}

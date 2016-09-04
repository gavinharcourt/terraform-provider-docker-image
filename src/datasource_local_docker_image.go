package dockerImage

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
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

			"tag": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A short human readable identifier for this image.",
				ForceNew:    true,
			},
		},
	}
}

func dataSourceLocalDockerImageCreate(d *schema.ResourceData, meta interface{}) error {
	pathToDockerfile := d.Get("dockerfile_path").(string)
	tag := d.Get("tag").(string)

	hash, err := dockerExec(meta.(*Config).DockerExecutable).buildContainer(pathToDockerfile, tag)
	if err != nil {
		return fmt.Errorf("Failed to create local docker image: %s", err)
	}

	d.SetId(hash)
	return nil
}

// following template data provider convention of doing real work in Exists.
// not sure how appropriate it is for a docker image.
func dataSourceLocalDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func dataSourceLocalDockerImageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	pathToDockerfile := d.Get("dockerfile_path").(string)
	tag := d.Get("tag").(string)

	hash, err := dockerExec(meta.(*Config).DockerExecutable).buildContainer(pathToDockerfile, tag)
	if err != nil {
		return false, fmt.Errorf("Failed to build local docker image: %s", err)
	}

	return hash == d.Id(), nil
}

func dataSourceLocalDockerImageDelete(d *schema.ResourceData, meta interface{}) error {
	err := dockerExec(meta.(*Config).DockerExecutable).deleteContainer(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to delete local docker image: %s", err)
	}
	return nil
}

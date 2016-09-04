package dockerImage

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceRemoteDockerImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceRemoteDockerImageCreate,
		Read:   resourceRemoteDockerImageRead,
		Update: resourceRemoteDockerImageUpdate,
		Delete: resourceRemoteDockerImageDelete,

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image ID to push to the remote.",
			},

			"tag": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tag of the remote docker image.",
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

func resourceRemoteDockerImageCreate(d *schema.ResourceData, meta interface{}) error {
	imageID := d.Get("image_id").(string)
	tag := d.Get("tag").(string)
	registry := d.Get("registry").(string)

	err := dockerExec(meta.(*Config).DockerExecutable).pushContainer(imageID, tag, registry)
	if err != nil {
		return fmt.Errorf("Failed to push container: %s", err)
	}

	return nil
}

func resourceRemoteDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("tag").(string))

	// set the image_id from the remote to "" so we always end up pushing
	d.Set("image_id", "")

	return nil
}

func resourceRemoteDockerImageUpdate(d *schema.ResourceData, meta interface{}) error {
	imageID := d.Get("image_id").(string)
	tag := d.Get("tag").(string)
	registry := d.Get("registry").(string)

	err := dockerExec(meta.(*Config).DockerExecutable).pushContainer(imageID, tag, registry)
	if err != nil {
		return fmt.Errorf("Failed to push container: %s", err)
	}

	return nil
}

func resourceRemoteDockerImageDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[WARN] deleting remote docker images is not currently supported, they must be deleted manually.")
	return nil
}

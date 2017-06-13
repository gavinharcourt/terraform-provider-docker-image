package main

import (
	"github.com/zongoose/terraform-provider-docker-image/src"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dockerImage.Provider,
	})
}

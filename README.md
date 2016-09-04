[![Build Status](https://travis-ci.org/diosmosis/terraform-provider-docker-image.svg?branch=master)](https://travis-ci.org/diosmosis/terraform-provider-docker-image)

# terraform-provider-docker-image

This provider contains two resource types that makes it easier to manage docker images & docker registries w/ terraform.

Currently, it's not possible to use terraform to push docker images to registries like ECR. If you want to push to a registry
as a build step, you'll have to do it via a bash script. This means there's no way to create a registry AND create ECS tasks
that depend on there being an image in that registry w/ one terraform project.

Instead, you have to create an ECR registry beforehand, push to it in a bash script, then build out the rest of your
infrastructure w/ terraform.

With this provider, you can get rid of this workflow, and use terraform to create your ECR registry & push an image to it.

## Resources

This plugin provides the following resources:

* `dockerimage_local` - references a local docker image by path. Every terraform apply will build the image to see if it has changed.
* `dockerimage_remote` - references a remote/local docker image mapping. Every terraform apply will attempt a push to the remote.
  _Note: it would be better if we only pushed when we knew the local changed, but it doesn't seem possible to check if a remote & local are the same._

## Setup

To use this plugin, you'll have to build it w/:

```
$ go build github.com/diosmosis/terraform-provider-docker-image
```

then add the following to your `~/.terraformrc` file:

```
providers {
  dockerimage = "/path/to/terraform-provider-docker-image"
}
```

## Example

```
provider "dockerimage" {}

resource "aws_ecr_repository" "myrepository" {
  name = "myrepository"
}

resource "dockerimage_local" "myimage" {
  dockerfile_path = "/path/to/dockerfiledirectory" # set to the path to the directory w/ your Dockerfile

  tag = "terraform-provider-docker-image-test:latest" # the tag for the image locally
}

resource "dockerimage_remote" "myimage" {
  tag = "terraform-provider-docker-image-test:latest" # the tag for the remote image
  registry = "${aws_ecr_repository.myrepository.registry_id}.dkr.ecr.us-east-1.amazonaws.com" # the registry's hostname
  image_id = "${dockerimage_local.myimage.id}" # the image ID to push
}
```

If using ECR, make sure you're logged in before running `terraform apply` (ie, w/ `$(aws ecr get-login --region us-east-1)`).

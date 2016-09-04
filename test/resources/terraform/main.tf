provider "dockerimage" {
  docker_executable = "docker"
}

resource "dockerimage_local" "test_image" {
  dockerfile_path = "${path.module}/.."
  tag = "terraform-provider-docker-image-test"
}

resource "dockerimage_remote" "test_image" {
  tag = "terraform-provider-docker-image-test:latest"
  registry = "${var.registry}"
  image_id = "${dockerimage_local.test_image.id}"
}

provider "dockerimage" {
  docker_executable = "docker"
}

provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region = "us-east-1"
}

resource "aws_ecr_repository" "test_repository" {
  name = "terraform-provider-docker-image-test"
}

resource "dockerimage_local" "test_image" {
  dockerfile_path = "${path.module}/.."
  tag = "terraform-provider-docker-image-test"
}

resource "dockerimage_remote" "test_image" {
  tag = "terraform-provider-docker-image-test:latest"
  registry = "${aws_ecr_repository.test_repository.registry_id}.dkr.ecr.us-east-1.amazonaws.com"
  image_id = "${dockerimage_local.test_image.id}"
}

#!/bin/bash

echo "providers {
  dockerImage = \"$TRAVIS_BUILD_DIR/terraform-provider-docker-image\"
}" > ~/.terraformrc

echo "registry = \"$ECR_REPOSITORY.dkr.ecr.us-east-1.amazonaws.com\"" &> test/resources/terraform/terraform.tfvars

export PATH="$PATH:$TRAVIS_BUILD_DIR"

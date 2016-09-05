#!/bin/bash

echo "providers {
  dockerimage = \"$TRAVIS_BUILD_DIR/terraform-provider-docker-image\"
}" > ~/.terraformrc

echo "registry = \"$ECR_REPOSITORY.dkr.ecr.us-east-1.amazonaws.com\"
aws_secret_key = \"$AWS_SECRET_ACCESS_KEY\"
aws_access_key = \"$AWS_ACCESS_KEY_ID\"
" &> test/resources/terraform/terraform.tfvars

export PATH="$PATH:$TRAVIS_BUILD_DIR"

cd test/resources/terraform

../../../terraform remote config -backend=s3 \
  -backend-config="bucket=terraform-provider-docker-image" \
  -backend-config="key=network/terraform.tfstate" \
  -backend-config="region=us-east-1"

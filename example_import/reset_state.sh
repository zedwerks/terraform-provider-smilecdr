#!/bin/bash
# Reset the state of the example_import project
rm terraform.tfstate
rm terraform.tfstate.backup
rm -rf .terraform.lock.hcl
terraform init -upgrade

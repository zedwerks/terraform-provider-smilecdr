#!/bin/sh
# script to update to a newer provider

rm .terraform.lock.hcl # remove lock file
terraform init -upgrade # upgrade to latest provider

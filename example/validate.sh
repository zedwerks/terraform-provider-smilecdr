#!/bin/sh
# script to update to a newer provider

terraform init -upgrade # upgrade to latest provider
terraform plan # check for errors
terraform validate # check for errors
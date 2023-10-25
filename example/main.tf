# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

terraform {
  required_providers {
    smilecdr = {
      source  = "local.providers/zedwerks/smilecdr"
      version = "~> 1.0.0"
    }
  }
}

provider "smilecdr" {
  base_url = "http://localhost:9000"
  username = "admin"
  password = "password"
}
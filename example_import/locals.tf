# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

locals {
  smilecdr_version             = "2021.05.R01"
  post_authorize_callback      = file("${path.module}/js/smart_outbound_security/post_authorize_callback.js")
}

# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

locals {
  smilecdr_version      = "2021.05.R01"
  callbackScriptFile    = file("${path.module}/js/inbound_security_callback.js")
  authScriptFile1       = file("${path.module}/js/federationAuthScript.js")
  authScriptFile2       = file("${path.module}/js/federationAuthScript.js")
  userMappingScriptFile = file("${path.module}/js/federationUserMappingScript.js")
}

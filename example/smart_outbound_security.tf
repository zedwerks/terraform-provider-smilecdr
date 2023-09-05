# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_smart_outbound_security" "smart_auth_outbound" {
  module_id = "smart_auth_outbound"
  node_id   = "Master"
  smart_capabilities_list = ""
}

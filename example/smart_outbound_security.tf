# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_smart_outbound_security" "smart_auth_outbound" {
  module_id                                       = "smart_auth_outbound"
  node_id                                         = "Master"
  smart_capabilities_list                         = ""
  pkce_required                                   = false
  pkce_plain_challenge_supported                  = false
  enforce_approved_scopes_to_restrict_permissions = false
  allowed_audience_list                           = ""
  federate_mode_enabled                           = false
  davinci_consent_handling_enabled                = false
}

# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0
#
# This is an example of how to use the smart inbound security module.
# An inbound security module validates tokens only. It is not an authorization server
# and does not issue tokens.  It is used to validate tokens issued by an external
# authorization server.



resource "smilecdr_smart_inbound_security" "inbound_security" {
    module_id = "smart_inbound"
    node_id = "Master"

}
// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/zed-werks/terraform-smilecdr/smilecdr"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return smilecdr.Provider()
		},
	})
}

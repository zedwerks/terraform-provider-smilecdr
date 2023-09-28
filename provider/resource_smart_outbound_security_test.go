// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSmartOutboundSecurity(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSmartOutboundConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccSmartOutboundModuleExists("smilecdr_smart_outbound_security.testacc"),
				),
			},
		},
	})
}

func testSmartOutboundConfig() string {

	moduleName := "smart_" + acctest.RandString(8)

	return fmt.Sprintf(`resource "smilecdr_smart_outbound_security" "testacc" {
		module_id                             = "%s"
		node_id                               = "Master"
		http_listener_context_path            = "/auth"
		http_listener_unhealthy_response_code = 503
		http_listener_port                    = 9999
		oidc_issuer_url                       = "http://keycloak:8080/auth/realms/lra"
		oidc_federate_mode_enabled            = true
		smart_authorization_enforce_approved_scopes = true
		smart_authorization_scopes_supported        = "launch fhirUser openid profile patient/*.read"
		sessions_max_concurrent_sessions_per_user   = 3
		dependency_fhir_persistence_module          = "PERSISTENCE_ALL"
}`, moduleName)
}

func testAccSmartOutboundModuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			fmt.Println("Not found: ", n)
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ClientId set")
		}

		return nil
	}
}

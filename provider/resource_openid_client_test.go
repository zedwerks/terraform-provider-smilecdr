// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func TestSmileCdrOpenIdClientBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccOpenIdClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMaxOpenIdClientConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccOpenIdClientExists("smilecdr_openid_client.max"),
				),
			},
			{
				ResourceName:      "smilecdr_openid_client.max",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMaxOpenIdClientConfig() string {
	return `resource "smilecdr_openid_client" "max" {		
		node_id                       = "Master"
		module_id                     = "smart_auth_federated"
		access_token_validity_seconds = 300
		allowed_grant_types           = ["REFRESH_TOKEN", "CLIENT_CREDENTIALS"]
		auto_approve_scopes = ["openid", "profile", "fhirUser", "launch", "launch/patient" ]
  		auto_grant_scopes = ["openid", "offline_access"]
  		client_id         = "max"
  		client_name       = "Client1"
  		client_secrets {
			secret     = "secret1234569900"
			activation = "2023-08-05T00:09:53.702+00:00"
		}
		client_secrets {
			secret = "secret23456789aabb"
		}
		fixed_scope                    = false
		refresh_token_validity_seconds = 86400
		registered_redirect_uris       = ["http://example-client1.com:6000", "http://example-client1.com:6001"]
		scopes = [
			"openid",
			"profile",
			"fhirUser",
			"patient/*.read",
			"launch",
			"launch/patient",
			"offline_access"
		]
		secret_required           = true
		secret_client_can_change  = false
		enabled                   = true
		can_introspect_any_tokens = true
		can_introspect_own_tokens = false
		always_require_approval   = false
		can_reissue_tokens        = false
		remember_approved_scopes  = false
		attestation_accepted      = false
		jwks_url                  = "http://example-client1.com/jwks"
		permissions {
			permission = "FHIR_WRITE_ALL_IN_COMPARTMENT"
			argument   = "Patient/123"
		}
		permissions {
			permission = "ROLE_FHIR_CLIENT"
		}	
	}`
}

func testAccOpenIdClientExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ClientId set")
		}

		return nil
	}
}

func testAccOpenIdClientDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*smilecdr.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "smilecdr_openid_client" {
			continue
		}

		nodeId := rs.Primary.Attributes["node_id"]
		moduleId := rs.Primary.Attributes["module_id"]
		clientId := rs.Primary.Attributes["client_id"]

		err := c.DeleteOpenIdClient(nodeId, moduleId, clientId)
		if err != nil {
			return err
		}
	}
	return nil
}

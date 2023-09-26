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

func TestSmileCdrOpenIdClientBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOpenIdClientConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccOpenIdClientExists("smilecdr_openid_client.basic"),
				),
			},
		},
	})
}

func TestSmileCdrOpenIdClientConfidential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOpenIdClientConfig_confidential(),
				Check: resource.ComposeTestCheckFunc(
					testAccOpenIdClientExists("smilecdr_openid_client.confidential"),
				),
			},
		},
	})
}

func testOpenIdClientConfig_basic() string {

	clientId := "cl_" + acctest.RandString(10)

	return fmt.Sprintf(`resource "smilecdr_openid_client" "basic" {		
		node_id                       = "Master"
		module_id                     = "smart_auth"
		access_token_validity_seconds = 300
		allowed_grant_types           = ["REFRESH_TOKEN", "CLIENT_CREDENTIALS"]
		auto_approve_scopes = ["openid", "profile", "fhirUser", "launch", "launch/patient" ]
  		auto_grant_scopes = ["openid", "offline_access"]
  		client_id         = "%s"
  		client_name       = "Client1"
		fixed_scope                    = false
		refresh_token_validity_seconds = 86400
		registered_redirect_uris       = ["http://example-client1.com:6000"]
		scopes = [
			"openid",
			"profile",
			"fhirUser",
			"patient/*.read",
			"launch",
			"offline_access"
		]
		secret_required           = false
		enabled                   = true
	}`, clientId)
}

func testOpenIdClientConfig_confidential() string {

	clientId := "cc_" + acctest.RandString(10)

	return fmt.Sprintf(`resource "smilecdr_openid_client" "confidential" {		
		node_id                       = "Master"
		module_id                     = "smart_auth"
		access_token_validity_seconds = 300
		allowed_grant_types           = ["REFRESH_TOKEN", "CLIENT_CREDENTIALS"]
		auto_approve_scopes = ["openid", "profile", "fhirUser", "launch", "launch/patient" ]
  		auto_grant_scopes = ["openid", "offline_access"]
  		client_id         = "%s"
  		client_name       = "Client2"
  		client_secrets {
			secret     = "secret1234569900"
			activation = "2023-08-05T00:09:53.702+00:00"
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

		# ignore secrets masking by smile cdr
		lifecycle {
			ignore_changes = [
				client_secrets,
			]
		}
	}`, clientId)
}

func testAccOpenIdClientExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			fmt.Println("Not found: ", n)
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			fmt.Println("No ClientId set")
			return fmt.Errorf("No ClientId set")
		}

		return nil
	}
}

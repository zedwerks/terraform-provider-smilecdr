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

func TestOpenIdProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOpenIdProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccOpenIdProviderExists("smilecdr_openid_identity_provider.testacc"),
				),
			},
		},
	})
}

func testOpenIdProviderConfig() string {

	providerName := "idp_" + acctest.RandString(10)

	return fmt.Sprintf(`resource "smilecdr_openid_identity_provider" "testacc" {		
		node_id                             = "Master"
		module_id                           = "smart_auth"
		issuer                              = "https://example1.com/auth/issuer"
		name                                = "%s"
		federation_request_scopes           = "openid profile email"
		federation_authorization_url        = "http://localhost:8800/auth"
		federation_token_url                = "http://localhost:8800/token"
		federation_user_info_url            = "http://localhost:8800/userinfo"
		federation_jwk_set_url              = "http://localhost:8800/auth/jwks"
		federation_auth_script_text         = "// This is a test script"
		token_introspection_client_id       = "smile"
		token_introspection_client_secret   = "client_secret_goes_here"
	}`, providerName)
}

func testAccOpenIdProviderExists(n string) resource.TestCheckFunc {
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

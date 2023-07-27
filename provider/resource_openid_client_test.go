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

	clientId := "test-client"
	clientName := "Test Client"
	callbackUrl := "https://example.com/callback"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmileCdrOpenIdClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSmileCdrOpenIdClientConfig(clientId, clientName, callbackUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmileCdrOpenIdClientExists("smilecdr_openid_client.new"),
				),
			},
		},
	})
}
func testAccCheckSmileCdrOpenIdClientConfig(clientId, name, callbackUrl string) string {
	return fmt.Sprintf(`
	resource "smilecdr_openid_client" "new" {
		{
			client_id = %s
    		client_name = %s
			allowed_grant_types = ["AUTHORIZATION_CODE"]
			callback_uris = [%s]
  		}
	}
	`, clientId, name, callbackUrl)
}

func testAccCheckSmileCdrOpenIdClientExists(n string) resource.TestCheckFunc {
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

func testAccCheckSmileCdrOpenIdClientDestroy(s *terraform.State) error {
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

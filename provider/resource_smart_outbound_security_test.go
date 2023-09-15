package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMyProvider_Resource(t *testing.T) {
	// Define the Terraform configuration for the test case
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartOutboundResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSmartOutboundResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartOutboundResourceExists("smilecdr_smart_outbound_security.example"),
					// Add more checks as needed
				),
			},
		},
	})
}

func testAccCheckSmartOutboundResourceConfig() string {
	return `resource "smilecdr_smart_outbound_security" "example" {
		node_id = "Master"
		module_id = "smart_outbound_security"
		fubar = true`
}

func testAccCheckSmartOutboundResourceDestroy(s *terraform.State) error {
	// Check if the resource was successfully destroyed
	// You may need to query your infrastructure to verify this
	return nil
}

func testAccCheckSmartOutboundResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Primary ID set")
		}

		return nil
	}
}

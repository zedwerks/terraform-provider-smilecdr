package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMyProvider_Resource(t *testing.T) {
	// Define the Terraform configuration for the test case
	resourceName := "smilecdr_smart_outbound_security.example"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartOutboundResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSmartOutboundResource_max(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartOutboundResourceExists(resourceName),
				),
			},
		},
	})
}

func testSmartOutboundResource_max() string {

	moduleId := acctest.RandString(11)

	return fmt.Sprintf(`resource "smilecdr_smart_outbound_security" "example" {
		module_id                             = "mid_%s"
		node_id                               = "Master"
		anonymous_account_username            = "ANONYMOUS"
		anonymous_access_enabled              = false
		saml_authentication_enabled           = false
		cors_allowed_headers                  = "Access-Control-Allow-Methods, Access-Control-Allow-Origin, Authorization, Content-Type, X-Requested-With"
		cors_enabled                          = true
		cors_origins                          = "*"
		davinci_native_consent_handling       = true
		http_access_log_appenders             = "appender-one\nappender-two"
		http_listener_context_path            = "/auth"
		http_listener_unhealthy_response_code = 503
		http_listener_port                    = 9201
		oidc_issuer_url                       = "http://keycloak:8080/auth/realms/poc"
		oidc_federate_mode_enabled            = false
		smart_authorization_enforce_approved_scopes = true
		smart_authorization_scopes_supported        = "launch fhirUser openid profile patient/*.read"
		sessions_max_concurrent_sessions_per_user   = 3
		dependency_fhir_persistence_module          = "PERSISTENCE_ALL"
	  }`, moduleId)
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

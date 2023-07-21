package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"smilecdr": testAccProvider,
	}
}

// Test the Provider
func TestProvider(t *testing.T) {
	// Create a new instance of the Provider
	provider := Provider()

	// Verify that the provider schema is valid
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("Error validating provider schema: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SMILECDR_BASE_URL"); v == "" {
		t.Fatal("SMILECDR_BASE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("SMILECDR_USERNAME"); v == "" {
		t.Fatal("SMILECDR_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("SMILECDR_PASSWORD"); v == "" {
		t.Fatal("SMILECDR_PASSWORD must be set for acceptance tests")
	}
}

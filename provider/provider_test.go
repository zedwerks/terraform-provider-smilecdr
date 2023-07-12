package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Test the Provider
func Test_Provider(t *testing.T) {
	// Create a new instance of the Provider
	provider := Provider()

	// Verify that the provider schema is valid
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("Error validating provider schema: %s", err)
	}
}

// Provider function
func Provider() *schema.Provider {
	return &schema.Provider{
		// Implement your provider's resources, data sources, and other functionality here
	}
}

func Test_envCheck(t *testing.T) {
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

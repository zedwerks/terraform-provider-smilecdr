// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestSmileCdrUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUserConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testUserExists("smilecdr_user.basic_user"),
				),
			},
		},
	})
}

func testUserConfig_basic() string {

	username := "U_" + strings.ToUpper(acctest.RandString(8))

	return fmt.Sprintf(`resource "smilecdr_user" "basic_user" {
		node_id = "Master"
		module_id = "local_security"
  		username = "%s"
		password = "Passw0rd"
		family_name = "Basic"
  		given_name = "User"
		account_locked = false
		system_user = false
		account_disabled = false
		external = false
		service_account = false

		authorities {
			permission = "FHIR_ALL_READ"
	    }
	    authorities {
		   permission = "ROLE_FHIR_CLIENT"
	    }

		# ignore password masking by smile cdr
		lifecycle {
			ignore_changes = [
				password,
			]
		}
	}`, username)
}

func testUserConfig_permissions() string {

	username := "U_" + strings.ToUpper(acctest.RandString(8))

	return fmt.Sprintf(`resource "smilecdr_user" "perm_user" {
		node_id = "Master"
		module_id = "local_security"
  		username = "%s"
		password = "Passw0rd"
		family_name = "Permissions"
  		given_name = "User"
		account_locked = false
		system_user = false
		account_disabled = false
		external = false
		service_account = false
		authorities {
			 { permission = "FHIR_ALL_READ"}
		}
		authorities {
			{ permission = "ROLE_FHIR_CLIENT"}
		}

		# ignore password masking by smile cdr
		lifecycle {
			ignore_changes = [
				password,
			]
		}
	}`, username)
}

func testUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			fmt.Println("Not found: ", n)
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			fmt.Println("No User set")
			return fmt.Errorf("No User set")
		}

		return nil
	}
}

package smilecdr

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func suppressSensitiveDataDiff(k, old, new string, d *schema.ResourceData) (bool, error) {
	// Your custom diff logic here
	log.Printf("Diff function called for key: %s, old: %s, new: %s\n", k, old, new)

	// For example, let's suppress the diff for the 'password' attribute
	if k == "password" {
		log.Println("Suppressing diff for sensitive attribute 'password'")
		return true, nil
	}
	if k == "secret" {
		log.Println("Suppressing diff for sensitive attribute 'secret'")
		return true, nil
	}

	// Continue with the default behavior for other attributes
	return false, nil
}

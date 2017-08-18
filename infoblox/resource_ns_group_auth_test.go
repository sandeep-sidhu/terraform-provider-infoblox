package infoblox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/skyinfoblox"
	"github.com/sky-uk/skyinfoblox/api/nsgroupauth"
	"regexp"
	"testing"
)

func TestAccInfobloxNSGroupAuthBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	nsGroupAuthNameCreate := fmt.Sprintf("acctest-infoblox-ns-group-auth-%d", randomInt)
	nsGroupAuthNameUpdate := fmt.Sprintf("%s-updated", nsGroupAuthNameCreate)
	nsGroupAuthResourceInstance := "infoblox_ns_group_auth.acctest"

	fmt.Printf("\n\nAcceptance Test NS Group Auth is %s\n\n", nsGroupAuthNameCreate)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccInfobloxNSGroupAuthCheckDestroy(state, nsGroupAuthNameCreate)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccInfobloxNSGroupAuthNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccInfobloxNSGroupAuthNameLeadingTrailingSpaces(),
				ExpectError: regexp.MustCompile(`must not contain trailing or leading white space`),
			},
			{
				Config:      testAccInfobloxNSGroupAuthCommentLeadingTrailingSpaces(),
				ExpectError: regexp.MustCompile(`must not contain trailing or leading white space`),
			},
			{
				Config: testAccInfobloxNSGroupAuthCreateTemplate(nsGroupAuthResourceInstance),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxNSGroupAuthCheckExists(nsGroupAuthNameCreate, nsGroupAuthResourceInstance),
					resource.TestCheckResourceAttr(nsGroupAuthResourceInstance, "name", nsGroupAuthNameCreate),
					resource.TestCheckResourceAttr(nsGroupAuthResourceInstance, "comment", "Infoblox Terraform Acceptance test"),
				),
			},
			{
				Config: testAccInfobloxNSGroupAuthUpdateTemplate(nsGroupAuthNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxAdminGroupCheckExists(nsGroupAuthNameUpdate, nsGroupAuthResourceInstance),
					resource.TestCheckResourceAttr(nsGroupAuthResourceInstance, "name", nsGroupAuthNameUpdate),
					resource.TestCheckResourceAttr(nsGroupAuthResourceInstance, "comment", "Infoblox Terraform Acceptance test - updated"),
				),
			},
		},
	})
}

func testAccInfobloxNSGroupAuthCheckDestroy(state *terraform.State, name string) error {

	client := testAccProvider.Meta().(*skyinfoblox.InfobloxClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_ns_group_auth" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		api := nsgroupauth.NewGetAll()
		err := client.Do(api)
		if err != nil {
			return fmt.Errorf("Infoblox - error occurred whilst retrieving a list of NS Group Auths")
		}
		for _, nsGroupAuth := range *api.ResponseObject().(*[]nsgroupauth.NSGroupAuth) {
			if nsGroupAuth.Name == name {
				return fmt.Errorf("Infoblox NS Group Auth %s still exists", name)
			}
		}
	}
	return nil
}

func testAccInfobloxNSGroupAuthCheckExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nInfoblox NS Group Auth %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nInfoblox NS Group Auth ID not set for %s in resources", name)
		}

		client := testAccProvider.Meta().(*skyinfoblox.InfobloxClient)
		api := nsgroupauth.NewGetAll()
		err := client.Do(api)
		if err != nil {
			return fmt.Errorf("Infoblox NS Group Auth - error whilst retrieving a list of NS Group Auths: %+v", err)
		}
		for _, nsGroupAuth := range *api.ResponseObject().(*[]nsgroupauth.NSGroupAuth) {
			if nsGroupAuth.Name == name {
				return nil
			}
		}
		return fmt.Errorf("Infoblox NS Group Auth %s wasn't found on remote Infoblox server", name)
	}
}

func testAccInfobloxNSGroupAuthNoNameTemplate() string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_auth" "acctest" {
comment = "Infoblox Terraform Acceptance test"
grid_default_group = true
use_external_primary = true
}
`)
}

func testAccInfobloxNSGroupAuthNameLeadingTrailingSpaces() string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_auth" "acctest" {
name = " test-group "
comment = "Infoblox Terraform Acceptance test"
grid_default_group = true
use_external_primary = true
}
`)
}

func testAccInfobloxNSGroupAuthCommentLeadingTrailingSpaces() string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_auth" "acctest" {
name = "test-group"
comment = " Infoblox Terraform Acceptance test "
grid_default_group = true
use_external_primary = true
}
`)
}

func testAccInfobloxNSGroupAuthCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_auth" "acctest" {
name = "%s"
comment = "Infoblox Terraform Acceptance test"
grid_default_group = true
use_external_primary = true
external_primaries = [
 {
     address = "192.168.0.1"
     name = "ns1.example.com"
     shared_with_ms_parent_delegation = false
     stealth = false
     tsig_key = "0jnu3SdsMvzzlmToPYRceA=="
     tsig_key_alg = "HMAC-MD5"
     tsig_key_name = "acc-test.key"
     use_tsig_key_name = true
 },
 ]
external_secondaries = [
  {
     address = "192.168.0.2"
     name = "ns2.example.com"
     shared_with_ms_parent_delegation = false
     stealth = false
     tsig_key = "0jnu3SdsMvzzlmToPYRceA=="
     tsig_key_alg = "HMAC-MD5"
     tsig_key_name = "acc-test.key"
     use_tsig_key_name = true
  },
]
grid_primary = [
  {
     gridreplicate = true
     lead = false
     name = "ns6.example.com"
     enablepreferredprimaries = true
     preferredprimaries = [
         {
           address = "192.168.1.1"
           name = "ns3.example.com"
           shared_with_ms_parent_delegation = false
           stealth = false
           tsig_key = "0jnu3SdsMvzzlmToPYRceA=="
           tsig_key_alg = "HMAC-MD5"
           tsig_key_name = "acc-test.key"
           use_tsig_key_name = true
         },
     ]
     stealth = false
  },
]
grid_secondaries = [
  {
     gridreplicate = true
     lead = false
     name = "ns5.example.com"
     enablepreferredprimaries = true
     preferredprimaries = [
         {
           address = "192.168.1.2"
           name = "ns4.example.com"
           shared_with_ms_parent_delegation = false
           stealth = false
           tsig_key = "0jnu3SdsMvzzlmToPYRceA=="
           tsig_key_alg = "HMAC-MD5"
           tsig_key_name = "acc-test.key"
           use_tsig_key_name = true
         },
     ]
     stealth = false
  },
]
}`, name)
}

func testAccInfobloxNSGroupAuthUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_auth" "acctest" {
name = "%s"
comment = "Infoblox Terraform Acceptance test - updated"
grid_default_group = false
use_external_primary = false
external_primaries = [
 {
     address = "192.168.10.1"
     name = "ns1.another-example.com"
     shared_with_ms_parent_delegation = false
     stealth = true
     tsig_key = "0jnu3SdsNvzzlmToPYRceA=="
     tsig_key_alg = "HMAC-SHA256"
     tsig_key_name = "acc-test2.key"
     use_tsig_key_name = true
 },
 ]
external_secondaries = [
  {
     address = "192.168.10.2"
     name = "ns2.another-example.com"
     shared_with_ms_parent_delegation = false
     stealth = true
     tsig_key = "0jnu3SdsNvzzlmToPYRceA=="
     tsig_key_alg = "HMAC-SHA256"
     tsig_key_name = "acc-test2.key"
     use_tsig_key_name = true
  },
]
grid_primary = [
  {
     gridreplicate = true
     lead = true
     name = "ns6.another-example.com"
     enablepreferredprimaries = true
     preferredprimaries = [
         {
           address = "192.168.10.1"
           name = "ns3.another-example.com"
           shared_with_ms_parent_delegation = false
           stealth = true
           tsig_key = "0jnu3SdsNvzzlmToPYRceA=="
           tsig_key_alg = "HMAC-sha256"
           tsig_key_name = "acc-test2.key"
           use_tsig_key_name = true
         },
     ]
     stealth = true
  },
]
grid_secondaries = [
  {
     gridreplicate = true
     lead = false
     name = "ns5.another-example.com"
     enablepreferredprimaries = true
     preferredprimaries = [
         {
           address = "192.168.10.2"
           name = "ns4.another-example.com"
           shared_with_ms_parent_delegation = false
           stealth = false
           tsig_key = "0jnu3SdsNvzzlmToPYRceA=="
           tsig_key_alg = "HMAC-sha256"
           tsig_key_name = "acc-test2.key"
           use_tsig_key_name = true
         },
     ]
     stealth = true
  },
]
}`, name)
}

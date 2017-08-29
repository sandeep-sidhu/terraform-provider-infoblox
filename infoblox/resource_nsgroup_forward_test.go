package infoblox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/skyinfoblox"
	"github.com/sky-uk/skyinfoblox/api/nsgroupfwd"
	"regexp"
	"testing"
)

func TestAccInfobloxNSGroupForwardBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	nsGroupForwardName := fmt.Sprintf("acctest-infoblox-ns-group-forward-%d", randomInt)
	nsGroupNameForwardUpdate := fmt.Sprintf("%s-updated", nsGroupForwardName)
	nsGroupForwardResourceInstance := "infoblox_ns_group_forward.acctest"

	fmt.Printf("\n\nAcceptance Test NS Group Forward is %s\n\n", nsGroupForwardName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccInfobloxNSGroupForwardCheckDestroy(state, nsGroupForwardName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccInfobloxNSGroupForwardNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccInfobloxNSGroupForwardCommentLeadingTrailingSpaces(nsGroupForwardName),
				ExpectError: regexp.MustCompile(`xmust not contain trailing or leading white space`),
			},
			{
				Config: testAccInfobloxNSGroupForwardCreateTemplate(nsGroupForwardName),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxNSGroupForwardCheckExists(nsGroupForwardName, nsGroupForwardResourceInstance),
					resource.TestCheckResourceAttr(nsGroupForwardResourceInstance, "name", nsGroupForwardName),
					resource.TestCheckResourceAttr(nsGroupForwardResourceInstance, "comment", "Infoblox Terraform Acceptance test"),
				),
			},
			{
				Config: testAccInfobloxNSGroupForwardUpdateTemplate(nsGroupNameForwardUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxNSGroupForwardCheckExists(nsGroupNameForwardUpdate, nsGroupForwardResourceInstance),
					resource.TestCheckResourceAttr(nsGroupForwardResourceInstance, "name", nsGroupNameForwardUpdate),
					resource.TestCheckResourceAttr(nsGroupForwardResourceInstance, "comment", "Infoblox Terraform Acceptance test - updated"),
				),
			},
		},
	})
}

func testAccInfobloxNSGroupForwardCheckDestroy(state *terraform.State, name string) error {

	client := testAccProvider.Meta().(*skyinfoblox.InfobloxClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "infoblox_ns_group_forward" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id == "" {
			return nil
		}
		api := nsgroupfwd.NewGetAll()
		err := client.Do(api)
		if err != nil {
			return fmt.Errorf("Infoblox - error occurred whilst retrieving a list of NS Group Foward")
		}
		for _, nsGroupForward := range *api.ResponseObject().(*[]nsgroupfwd.NSGroupFwd) {
			if nsGroupForward.Name == name {
				return fmt.Errorf("Infoblox NS Group Forward %s still exists", name)
			}
		}
	}
	return nil
}

func testAccInfobloxNSGroupForwardCheckExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("\nInfoblox NS Group Forward %s wasn't found in resources", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nInfoblox NS Group Forward ID not set for %s in resources", name)
		}

		client := testAccProvider.Meta().(*skyinfoblox.InfobloxClient)
		api := nsgroupfwd.NewGetAll()
		err := client.Do(api)
		if err != nil {
			return fmt.Errorf("Infoblox NS Group Forward - error whilst retrieving a list of NS Group Forward: %+v", err)
		}
		for _, nsGroupForward := range *api.ResponseObject().(*[]nsgroupfwd.NSGroupFwd) {
			if nsGroupForward.Name == name {
				return nil
			}
		}
		return fmt.Errorf("Infoblox NS Group Forward %s wasn't found on remote Infoblox server", name)
	}
}

func testAccInfobloxNSGroupForwardNoNameTemplate() string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_forward" "acctest" {
  comment = "Infoblox Terraform Acceptance test"
}
`)
}

func testAccInfobloxNSGroupForwardCommentLeadingTrailingSpaces(name string) string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_forward" "acctest" {
  name = "%s"
  comment = " Infoblox Terraform Acceptance test "
}
`, name)
}

func testAccInfobloxNSGroupForwardCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_forward" "acctest" {
  name = "%s"
  comment = "Infoblox Terraform Acceptance test"
  forwarding_servers = [
    {
      forward_to = [
        {
	  name = "ns1.example.com"
	  address = "192.168.1.3"
        },
        {
	  name = "ns2.example.com"
	  address = "192.168.1.4"
        },
      ],
      forwarders_only = true
      name = ""
      use_override_forwarders = true
    },
    {
      forward_to = [
        {
	  name = "ns11.example.com"
	  address = "192.168.2.3"
        },
        {
	  name = "ns12.example.com"
	  address = "192.168.2.4"
        },
      ],
      forwarders_only = true
      name = ""
      use_override_forwarders = true
    },
  ],
}
`, name)
}

func testAccInfobloxNSGroupForwardUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "infoblox_ns_group_forward" "acctest" {
  name = "%s"
  comment = "Infoblox Terraform Acceptance test - updated"
}
`, name)
}

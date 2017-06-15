package infoblox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/skyinfoblox"
	"github.com/sky-uk/skyinfoblox/api/records"
	"regexp"
	"testing"
)

func TestAccInfobloxCNAMEBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	cname := fmt.Sprintf("acctest-infoblox-cname-%d.ovp.bskyb.com", randomInt)
	canonical := fmt.Sprintf("acctest-infoblox-canonical-%d.ovp.bskyb.com", randomInt)
	canonicalUpdate := fmt.Sprintf("acctest-infoblox-canonical-update-%d.ovp.bskyb.com", randomInt)
	cnameResourceName := "infoblox_cname_record.acctest"

	fmt.Printf("\n\nCNAME is %s\n\n", cname)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccInfobloxCNAMECheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInfobloxCNAMECreateTemplate(cname, canonical),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxCNAMEExists(cname, cnameResourceName),
					resource.TestCheckResourceAttr(cnameResourceName, "name", cname),
					resource.TestCheckResourceAttr(cnameResourceName, "comment", "Terraform Acceptance Testing for CNAMEs"),
					resource.TestCheckResourceAttr(cnameResourceName, "canonical", canonical),
					resource.TestCheckResourceAttr(cnameResourceName, "view", "default"),
					resource.TestCheckResourceAttr(cnameResourceName, "ttl", "1202"),
				),
			},
			{
				Config: testAccInfobloxCNAMEUpdateTemplate(cname, canonicalUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxCNAMEExists(cname, cnameResourceName),
					resource.TestCheckResourceAttr(cnameResourceName, "name", cname),
					resource.TestCheckResourceAttr(cnameResourceName, "comment", "Terraform Acceptance Testing for CNAMEs update test"),
					resource.TestCheckResourceAttr(cnameResourceName, "canonical", canonicalUpdate),
					resource.TestCheckResourceAttr(cnameResourceName, "view", "default"),
					resource.TestCheckResourceAttr(cnameResourceName, "ttl", "600"),
				),
			},
			/*
				{
					Config:      testAccInfobloxCNAMEInvalidTTLTemplate(cname, canonical),
					ExpectError: regexp.MustCompile(`.*`),

				},
			*/
		},
	})
}

func testAccInfobloxCNAMEExists(cnameCheck, cnameResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		returnFields := []string{"name", "comment", "view", "ttl", "canonical", "zone"}

		rs, ok := state.RootModule().Resources[cnameResourceName]
		if !ok {
			return fmt.Errorf("Infoblox CNAME resource %s not found in resources", cnameResourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Infoblox CNAME resource ID not set in resources")
		}

		client := testAccProvider.Meta().(*skyinfoblox.InfobloxClient)
		getAllAPI := records.NewGetAllCNAMERecords(returnFields)

		err := client.Do(getAllAPI)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, cname := range getAllAPI.GetResponse() {

			if cname.Name == cnameCheck {
				return nil
			}
		}
		return fmt.Errorf("Infoblox CNAME %s wasn't found", cnameCheck)
	}
}

func testAccInfobloxCNAMECheckDestroy(state *terraform.State) error {

	infobloxClient := testAccProvider.Meta().(*skyinfoblox.InfobloxClient)
	returnFields := []string{"name", "comment", "view", "ttl", "canonical", "zone"}

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "infoblox_cname_record" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		api := records.NewGetAllCNAMERecords(returnFields)
		err := infobloxClient.Do(api)
		if err != nil {
			return nil
		}
		for _, cname := range api.GetResponse() {
			matched, _ := regexp.MatchString("acctest-infoblox-cname-.*.ovp.bskyb.com", cname.Name)
			if matched {
				return fmt.Errorf("Sky Infoblox CNAME %s still exists", cname.Name)
			}
		}
	}
	return nil
}

func testAccInfobloxCNAMECreateTemplate(cname, canonical string) string {
	return fmt.Sprintf(`
resource "infoblox_cname_record" "acctest" {
  name = "%s"
  comment = "Terraform Acceptance Testing for CNAMEs"
  canonical = "%s"
  view = "default"
  ttl = 1202
}
`, cname, canonical)
}

func testAccInfobloxCNAMEUpdateTemplate(cname, canonical string) string {
	return fmt.Sprintf(`
resource "infoblox_cname_record" "acctest" {
  name = "%s"
  comment = "Terraform Acceptance Testing for CNAMEs update test"
  canonical = "%s"
  view = "default"
  ttl = 600
}
`, cname, canonical)
}

func testAccInfobloxCNAMEInvalidTTLTemplate(cname, canonical string) string {
	return fmt.Sprintf(`
resource "infoblox_cname_record" "acctest" {
  name = "%s"
  comment = "Terraform Acceptance Testing for CNAMEs update test"
  canonical = "%s"
  view = "default"
  ttl = -600
}
`, cname, canonical)
}

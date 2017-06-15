package infoblox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/skyinfoblox"
	"github.com/sky-uk/skyinfoblox/api/records"
	"strings"
)

func resourceCNAMERecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceCNAMECreate,
		Read:   resourceCNAMERead,
		Update: resourceCNAMEUpdate,
		Delete: resourceCNAMEDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"view": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			//"zone": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				// Implement validator function unsigned int.
			},
			"canonical": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCNAMECreate(d *schema.ResourceData, m interface{}) error {

	infobloxClient := m.(*skyinfoblox.InfobloxClient)
	var cnameRecord records.GenericRecord
	recordType := "cname"

	if v, ok := d.GetOk("name"); ok {
		cnameRecord.Name = v.(string)
	} else {
		return fmt.Errorf("Sky Infoblox Create Error: name argument required")
	}
	if v, ok := d.GetOk("comment"); ok {
		cnameRecord.Comment = v.(string)
	}
	if v, ok := d.GetOk("view"); ok {
		cnameRecord.View = v.(string)
	}
	//if v, ok := d.GetOk("zone"); ok {
	//	cnameRecord.Zone = v.(string)
	//}
	if v, ok := d.GetOk("ttl"); ok {
		ttl := v.(int)
		cnameRecord.TTL = uint(ttl)
	}
	if v, ok := d.GetOk("canonical"); ok {
		cnameRecord.Canonical = v.(string)
	}

	createAPI := records.NewCreateRecord(recordType, cnameRecord)

	err := infobloxClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Sky Infoblox Create Error: %+v", err)
	}

	if createAPI.StatusCode() != 201 {
		return fmt.Errorf("Sky Infoblox Create Error: Invalid HTTP response code %+v returned. Response object was %+v", createAPI.StatusCode(), createAPI.ResponseObject())
	}

	id := strings.Replace(createAPI.GetResponse(), "\"", "", -1)
	d.SetId(id)
	return resourceCNAMERead(d, m)
}

func resourceCNAMERead(d *schema.ResourceData, m interface{}) error {

	returnFields := []string{"name", "comment", "view", "ttl", "canonical"}

	infobloxClient := m.(*skyinfoblox.InfobloxClient)
	getSingleCNAMEAPI := records.NewGetCNAMERecord(d.Id(), returnFields)

	err := infobloxClient.Do(getSingleCNAMEAPI)
	if err != nil {
		return fmt.Errorf("Sky Infoblox Read Error: %+v", err)
	}
	if getSingleCNAMEAPI.StatusCode() == 404 {
		d.SetId("")
		return nil
	}

	response := getSingleCNAMEAPI.GetResponse()
	d.SetId(response.Ref)
	d.Set("name", response.Name)
	d.Set("comment", response.Comment)
	d.Set("view", response.View)
	d.Set("ttl", response.TTL)
	d.Set("canonical", response.Canonical)

	return nil
}

func resourceCNAMEUpdate(d *schema.ResourceData, m interface{}) error {

	return resourceCNAMERead(d, m)
}

func resourceCNAMEDelete(d *schema.ResourceData, m interface{}) error {

	d.SetId("")
	return nil
}
